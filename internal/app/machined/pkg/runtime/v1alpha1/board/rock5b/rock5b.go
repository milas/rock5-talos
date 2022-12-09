// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package rock5b

import (
	"log"
	"os"
	"path/filepath"

	"github.com/siderolabs/go-procfs/procfs"
	"golang.org/x/sys/unix"

	"github.com/siderolabs/talos/internal/app/machined/pkg/runtime"
	"github.com/siderolabs/talos/pkg/copy"
	"github.com/siderolabs/talos/pkg/machinery/constants"
)

const (
	sectorSize        = 512
	ubootImage        = "/usr/install/arm64/u-boot/rock_5b/u-boot.img"
	ubootOffset int64 = sectorSize * 0x40
	dtb               = "/dtb/rockchip/rk3588-rock-5b.dtb"
)

// Rock5b represents the Radxa Rock 5B board.
//
// Reference: https://wiki.radxa.com/Rock5
type Rock5b struct{}

// Name implements the runtime.Board.
func (r *Rock5b) Name() string {
	return constants.BoardRock5b
}

// Install implements the runtime.Board.
func (r *Rock5b) Install(disk string) (err error) {
	var f *os.File

	if f, err = os.OpenFile(disk, os.O_RDWR|unix.O_CLOEXEC, 0o666); err != nil {
		return err
	}

	defer f.Close() //nolint:errcheck

	uboot, err := os.ReadFile(ubootImage)
	if err != nil {
		return err
	}
	uboot = uboot[ubootOffset:]

	log.Printf("writing %s (%d) at offset %d", ubootImage, len(uboot), ubootOffset)

	var n int

	n, err = f.WriteAt(uboot, ubootOffset)
	if err != nil {
		return err
	}

	log.Printf("wrote %d bytes", n)

	// NB: In the case that the block device is a loopback device, we sync here
	// to ensure that the file is written before the loopback device is
	// unmounted.
	err = f.Sync()
	if err != nil {
		return err
	}

	src := "/usr/install/arm64" + dtb
	dst := "/boot/EFI" + dtb

	err = os.MkdirAll(filepath.Dir(dst), 0o600)
	if err != nil {
		return err
	}

	return copy.File(src, dst)
}

// KernelArgs implements the runtime.Board.
func (r *Rock5b) KernelArgs() procfs.Parameters {
	return []*procfs.Parameter{
		procfs.NewParameter("console").Append("tty0").Append("ttyS2,1500000n8"),
		procfs.NewParameter("sysctl.kernel.kexec_load_disabled").Append("1"),
		procfs.NewParameter(constants.KernelParamDashboardDisabled).Append("1"),
	}
}

// PartitionOptions implements the runtime.Board.
func (r *Rock5b) PartitionOptions() *runtime.PartitionOptions {
	return nil
}
