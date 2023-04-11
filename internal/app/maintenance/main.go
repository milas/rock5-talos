// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package maintenance

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/netip"
	"time"

	"github.com/cosi-project/runtime/pkg/resource"
	"github.com/cosi-project/runtime/pkg/state"
	ttls "github.com/siderolabs/crypto/tls"
	"github.com/siderolabs/crypto/x509"
	"github.com/siderolabs/gen/slices"
	"github.com/siderolabs/gen/value"
	"github.com/siderolabs/go-procfs/procfs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/siderolabs/talos/internal/app/machined/pkg/runtime"
	"github.com/siderolabs/talos/internal/app/maintenance/server"
	"github.com/siderolabs/talos/pkg/grpc/factory"
	"github.com/siderolabs/talos/pkg/grpc/gen"
	"github.com/siderolabs/talos/pkg/grpc/middleware/authz"
	"github.com/siderolabs/talos/pkg/machinery/constants"
	"github.com/siderolabs/talos/pkg/machinery/resources/network"
)

var ctrl runtime.Controller

// InjectController is used to pass the controller into the maintenance service.
func InjectController(c runtime.Controller) {
	ctrl = c
}

// Run executes the configuration receiver, returning any configuration it receives.
//
//nolint:gocyclo
func Run(ctx context.Context, logger *log.Logger) ([]byte, error) {
	if ctrl == nil {
		return nil, fmt.Errorf("controller is not injected")
	}

	logger.Println("waiting for network address to be ready")

	if err := network.NewReadyCondition(ctrl.Runtime().State().V1Alpha2().Resources(), network.AddressReady).Wait(ctx); err != nil {
		return nil, fmt.Errorf("error waiting for the network to be ready: %w", err)
	}

	logger.Println("loading current addresses")

	var sideroLinkAddress netip.Addr

	currentAddresses, err := ctrl.Runtime().State().V1Alpha2().Resources().WatchFor(
		ctx,
		resource.NewMetadata(
			network.NamespaceName,
			network.NodeAddressType,
			network.NodeAddressCurrentID,
			resource.VersionUndefined,
		),
		sideroLinkAddressFinder(&sideroLinkAddress, logger),
	)
	if err != nil {
		return nil, fmt.Errorf("error getting node addresses: %w", err)
	}

	ips := currentAddresses.(*network.NodeAddress).TypedSpec().IPs()

	// hostname might not be available yet, so use it only if it is available
	hostnameStatus, err := ctrl.Runtime().State().V1Alpha2().Resources().Get(ctx, resource.NewMetadata(network.NamespaceName, network.HostnameStatusType, network.HostnameID, resource.VersionUndefined))
	if err != nil && !state.IsNotFoundError(err) {
		return nil, fmt.Errorf("error getting node hostname: %w", err)
	}

	var dnsNames []string

	if hostnameStatus != nil {
		dnsNames = hostnameStatus.(*network.HostnameStatus).TypedSpec().DNSNames()
	}

	logger.Println("generating TLS config")

	tlsConfig, provider, err := genTLSConfig(ips, dnsNames)
	if err != nil {
		return nil, err
	}

	cert, err := provider.GetCertificate(nil)
	if err != nil {
		return nil, err
	}

	logger.Println("fingerprinting cert")

	certFingerprint, err := x509.SPKIFingerprintFromDER(cert.Certificate[0])
	if err != nil {
		return nil, err
	}

	cfgCh := make(chan []byte)

	logger.Println("* creating server")

	s := server.New(ctrl, logger, cfgCh)

	injector := &authz.Injector{
		Mode:   authz.ReadOnly,
		Logger: log.New(logger.Writer(), "machined/authz/injector ", log.Flags()).Printf,
	}

	// Start the server.
	server := factory.NewServer(
		s,
		factory.WithDefaultLog(),
		factory.ServerOptions(
			grpc.Creds(
				credentials.NewTLS(tlsConfig),
			),
		),

		factory.WithUnaryInterceptor(injector.UnaryInterceptor()),
		factory.WithStreamInterceptor(injector.StreamInterceptor()),
	)

	logger.Println("creating listener")

	listener, err := factory.NewListener(factory.Address(formatIP(sideroLinkAddress)), factory.Port(constants.ApidPort))
	if err != nil {
		return nil, err
	}

	defer func() {
		shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 5*time.Second)
		defer shutdownCancel()

		factory.ServerGracefulStop(server, shutdownCtx)
	}()

	go func() {
		//nolint:errcheck
		server.Serve(listener)
	}()

	if !value.IsZero(sideroLinkAddress) {
		ips = []netip.Addr{sideroLinkAddress}
	}

	logger.Println("this machine is reachable at:")

	for _, ip := range ips {
		logger.Printf("\t%s", ip.String())
	}

	firstIP := "<IP>"

	if len(ips) > 0 {
		firstIP = ips[0].String()
	}

	logger.Println("server certificate fingerprint:")
	logger.Printf("\t%s", certFingerprint)

	logger.Println()
	logger.Println("upload configuration using talosctl:")
	logger.Printf("\ttalosctl apply-config --insecure --nodes %s --file <config.yaml>", firstIP)
	logger.Println("or apply configuration using talosctl interactive installer:")
	logger.Printf("\ttalosctl apply-config --insecure --nodes %s --mode=interactive", firstIP)
	logger.Println("optionally with node fingerprint check:")
	logger.Printf(
		"\ttalosctl apply-config --insecure --nodes %s --cert-fingerprint '%s' --file <config.yaml>",
		firstIP,
		certFingerprint,
	)

	select {
	case cfg := <-cfgCh:
		return cfg, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func formatIP(addr netip.Addr) string {
	if value.IsZero(addr) {
		return ""
	}

	return addr.String()
}

func genTLSConfig(
	ips []netip.Addr,
	dnsNames []string,
) (tlsConfig *tls.Config, provider ttls.CertificateProvider, err error) {
	ca, err := x509.NewSelfSignedCertificateAuthority()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate self-signed CA: %w", err)
	}

	ips = append(ips, netip.MustParseAddr("127.0.0.1"), netip.MustParseAddr("::1"))

	netIPs := slices.Map(ips, func(ip netip.Addr) net.IP { return ip.AsSlice() })

	var generator ttls.Generator

	generator, err = gen.NewLocalGenerator(ca.KeyPEM, ca.CrtPEM)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create local generator provider: %w", err)
	}

	provider, err = ttls.NewRenewingCertificateProvider(generator, x509.DNSNames(dnsNames), x509.IPAddresses(netIPs))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create local certificate provider: %w", err)
	}

	caCertPEM, err := provider.GetCA()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get CA: %w", err)
	}

	tlsConfig, err = ttls.New(
		ttls.WithClientAuthType(ttls.ServerOnly),
		ttls.WithCACertPEM(caCertPEM),
		ttls.WithServerCertificateProvider(provider),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate tlsconfig: %w", err)
	}

	return tlsConfig, provider, nil
}

func sideroLinkAddressFinder(address *netip.Addr, logger *log.Logger) state.WatchForConditionFunc {
	sideroLinkEnabled := false
	if procfs.ProcCmdline().Get(constants.KernelParamSideroLink).First() != nil {
		sideroLinkEnabled = true

		logger.Println(constants.KernelParamSideroLink + " is enabled, waiting for address")
	}

	return state.WithCondition(
		func(r resource.Resource) (bool, error) {
			if resource.IsTombstone(r) {
				return false, nil
			}

			if !sideroLinkEnabled {
				return true, nil
			}

			ips := r.(*network.NodeAddress).TypedSpec().IPs()
			for _, ip := range ips {
				if network.IsULA(ip, network.ULASideroLink) {
					*address = ip

					return true, nil
				}
			}

			return false, nil
		},
	)
}
