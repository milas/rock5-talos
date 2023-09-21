// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package imager

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
	"github.com/google/go-containerregistry/pkg/v1/types"
	"github.com/siderolabs/go-procfs/procfs"

	"github.com/siderolabs/talos/cmd/installer/pkg/install"
	"github.com/siderolabs/talos/internal/app/machined/pkg/runtime/v1alpha1/bootloader/options"
	"github.com/siderolabs/talos/pkg/imager/filemap"
	"github.com/siderolabs/talos/pkg/imager/iso"
	"github.com/siderolabs/talos/pkg/imager/ova"
	"github.com/siderolabs/talos/pkg/imager/profile"
	"github.com/siderolabs/talos/pkg/imager/qemuimg"
	"github.com/siderolabs/talos/pkg/imager/utils"
	"github.com/siderolabs/talos/pkg/machinery/constants"
	"github.com/siderolabs/talos/pkg/reporter"
)

func (i *Imager) outInitramfs(path string, report *reporter.Reporter) error {
	printf := progressPrintf(report, reporter.Update{Message: "copying initramfs...", Status: reporter.StatusRunning})

	if err := utils.CopyFiles(printf, utils.SourceDestination(i.initramfsPath, path)); err != nil {
		return err
	}

	report.Report(reporter.Update{Message: "initramfs output ready", Status: reporter.StatusSucceeded})

	return nil
}

func (i *Imager) outKernel(path string, report *reporter.Reporter) error {
	printf := progressPrintf(report, reporter.Update{Message: "copying kernel...", Status: reporter.StatusRunning})

	if err := utils.CopyFiles(printf, utils.SourceDestination(i.prof.Input.Kernel.Path, path)); err != nil {
		return err
	}

	report.Report(reporter.Update{Message: "kernel output ready", Status: reporter.StatusSucceeded})

	return nil
}

func (i *Imager) outUKI(path string, report *reporter.Reporter) error {
	printf := progressPrintf(report, reporter.Update{Message: "copying kernel...", Status: reporter.StatusRunning})

	if err := utils.CopyFiles(printf, utils.SourceDestination(i.ukiPath, path)); err != nil {
		return err
	}

	report.Report(reporter.Update{Message: "UKI output ready", Status: reporter.StatusSucceeded})

	return nil
}

func (i *Imager) outISO(path string, report *reporter.Reporter) error {
	printf := progressPrintf(report, reporter.Update{Message: "building ISO...", Status: reporter.StatusRunning})

	scratchSpace := filepath.Join(i.tempDir, "iso")

	var err error

	if i.prof.SecureBootEnabled() {
		err = iso.CreateUEFI(printf, iso.UEFIOptions{
			UKIPath:    i.ukiPath,
			SDBootPath: i.sdBootPath,

			PlatformKeyPath:    i.prof.Input.SecureBoot.PlatformKeyPath,
			KeyExchangeKeyPath: i.prof.Input.SecureBoot.KeyExchangeKeyPath,
			SignatureKeyPath:   i.prof.Input.SecureBoot.SignatureKeyPath,

			Arch:    i.prof.Arch,
			Version: i.prof.Version,

			ScratchDir: scratchSpace,
			OutPath:    path,
		})
	} else {
		err = iso.CreateGRUB(printf, iso.GRUBOptions{
			KernelPath:    i.prof.Input.Kernel.Path,
			InitramfsPath: i.initramfsPath,
			Cmdline:       i.cmdline,

			ScratchDir: scratchSpace,
			OutPath:    path,
		})
	}

	if err != nil {
		return err
	}

	report.Report(reporter.Update{Message: "ISO ready", Status: reporter.StatusSucceeded})

	return nil
}

func (i *Imager) outImage(ctx context.Context, path string, report *reporter.Reporter) error {
	printf := progressPrintf(report, reporter.Update{Message: "creating disk image...", Status: reporter.StatusRunning})

	if err := i.buildImage(ctx, path, printf); err != nil {
		return err
	}

	switch i.prof.Output.ImageOptions.DiskFormat {
	case profile.DiskFormatRaw:
		// nothing to do
	case profile.DiskFormatQCOW2:
		if err := qemuimg.Convert("raw", "qcow2", i.prof.Output.ImageOptions.DiskFormatOptions, path, printf); err != nil {
			return err
		}
	case profile.DiskFormatVPC:
		if err := qemuimg.Convert("raw", "vpc", i.prof.Output.ImageOptions.DiskFormatOptions, path, printf); err != nil {
			return err
		}
	case profile.DiskFormatOVA:
		scratchPath := filepath.Join(i.tempDir, "ova")

		if err := ova.CreateOVAFromRAW(path, i.prof.Arch, scratchPath, i.prof.Output.ImageOptions.DiskSize, printf); err != nil {
			return err
		}
	case profile.DiskFormatUnknown:
		fallthrough
	default:
		return fmt.Errorf("unsupported disk format: %s", i.prof.Output.ImageOptions.DiskFormat)
	}

	report.Report(reporter.Update{Message: "disk image ready", Status: reporter.StatusSucceeded})

	return nil
}

func (i *Imager) buildImage(ctx context.Context, path string, printf func(string, ...any)) error {
	if err := utils.CreateRawDisk(printf, path, i.prof.Output.ImageOptions.DiskSize); err != nil {
		return err
	}

	printf("attaching loopback device")

	var (
		loDevice string
		err      error
	)

	if loDevice, err = utils.Loattach(path); err != nil {
		return err
	}

	defer func() {
		printf("detaching loopback device")

		if e := utils.Lodetach(loDevice); e != nil {
			log.Println(e)
		}
	}()

	cmdline := procfs.NewCmdline(i.cmdline)

	opts := &install.Options{
		Disk:       loDevice,
		Platform:   i.prof.Platform,
		Arch:       i.prof.Arch,
		Board:      i.prof.Board,
		MetaValues: install.FromMeta(i.prof.Customization.MetaContents),

		ImageSecureboot: i.prof.SecureBootEnabled(),
		Version:         i.prof.Version,
		BootAssets: options.BootAssets{
			KernelPath:    i.prof.Input.Kernel.Path,
			InitramfsPath: i.initramfsPath,
			UKIPath:       i.ukiPath,
			SDBootPath:    i.sdBootPath,
		},
		Printf: printf,
	}

	if !strings.HasPrefix(opts.Board, "rock") {
		panic("rock5: imager only works for rock5 rock5")
	}

	installer, err := install.NewInstaller(ctx, cmdline, install.ModeImage, opts)
	if err != nil {
		return fmt.Errorf("failed to create installer: %w", err)
	}

	if err := installer.Install(ctx, install.ModeImage); err != nil {
		return fmt.Errorf("failed to install: %w", err)
	}

	return nil
}

//nolint:gocyclo
func (i *Imager) outInstaller(ctx context.Context, path string, report *reporter.Reporter) error {
	printf := progressPrintf(report, reporter.Update{Message: "building installer...", Status: reporter.StatusRunning})

	baseInstallerImg, err := i.prof.Input.BaseInstaller.Pull(ctx, i.prof.Arch, printf)
	if err != nil {
		return err
	}

	baseLayers, err := baseInstallerImg.Layers()
	if err != nil {
		return fmt.Errorf("failed to get layers: %w", err)
	}

	configFile, err := baseInstallerImg.ConfigFile()
	if err != nil {
		return fmt.Errorf("failed to get config file: %w", err)
	}

	config := *configFile.Config.DeepCopy()

	printf("creating empty image")

	newInstallerImg := mutate.MediaType(empty.Image, types.OCIManifestSchema1)
	newInstallerImg = mutate.ConfigMediaType(newInstallerImg, types.OCIConfigJSON)

	newInstallerImg, err = mutate.Config(newInstallerImg, config)
	if err != nil {
		return fmt.Errorf("failed to set config: %w", err)
	}

	newInstallerImg, err = mutate.CreatedAt(newInstallerImg, v1.Time{Time: time.Now()})
	if err != nil {
		return fmt.Errorf("failed to set created at: %w", err)
	}

	newInstallerImg, err = mutate.AppendLayers(newInstallerImg, baseLayers[0])
	if err != nil {
		return fmt.Errorf("failed to append layers: %w", err)
	}

	var artifacts []filemap.File

	printf("generating artifacts layer")

	if i.prof.SecureBootEnabled() {
		artifacts = append(artifacts,
			filemap.File{
				ImagePath:  fmt.Sprintf(constants.UKIAssetPath, i.prof.Arch),
				SourcePath: i.ukiPath,
			},
			filemap.File{
				ImagePath:  fmt.Sprintf(constants.SDBootAssetPath, i.prof.Arch),
				SourcePath: i.sdBootPath,
			},
		)
	} else {
		artifacts = append(artifacts,
			filemap.File{
				ImagePath:  fmt.Sprintf(constants.KernelAssetPath, i.prof.Arch),
				SourcePath: i.prof.Input.Kernel.Path,
			},
			filemap.File{
				ImagePath:  fmt.Sprintf(constants.InitramfsAssetPath, i.prof.Arch),
				SourcePath: i.initramfsPath,
			},
		)
	}

	artifactsLayer, err := filemap.Layer(artifacts)
	if err != nil {
		return fmt.Errorf("failed to create artifacts layer: %w", err)
	}

	newInstallerImg, err = mutate.AppendLayers(newInstallerImg, artifactsLayer)
	if err != nil {
		return fmt.Errorf("failed to append artifacts layer: %w", err)
	}

	ref, err := name.ParseReference(i.prof.Input.BaseInstaller.ImageRef)
	if err != nil {
		return fmt.Errorf("failed to parse image reference: %w", err)
	}

	printf("writing image tarball")

	if err := tarball.WriteToFile(path, ref, newInstallerImg); err != nil {
		return fmt.Errorf("failed to write image tarball: %w", err)
	}

	report.Report(reporter.Update{Message: "installer container image ready", Status: reporter.StatusSucceeded})

	return nil
}
