// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package profile

import (
	"fmt"
	"github.com/siderolabs/go-pointer"
	"path/filepath"

	"github.com/siderolabs/talos/pkg/machinery/constants"
)

const (
	mib = 1024 * 1024

	// MinRAWDiskSize is the minimum size disk we can create. Used for metal images.
	MinRAWDiskSize = 1246 * mib

	// DefaultRAWDiskSize is the value we use for any non-metal images by default.
	DefaultRAWDiskSize = 8192 * mib
)

// Default describes built-in profiles.
var Default = map[string]Profile{
	// SBCs
	constants.BoardRock5a: {
		Arch:       "arm64",
		Platform:   constants.PlatformMetal,
		Board: constants.BoardRock5a,
		SecureBoot: pointer.To(false),
		Input: Input{
			DtbPath: FileAsset{Path: filepath.Join(fmt.Sprintf(constants.DtbsAssetPath, "arm64"), "rockchip", "rk3588s-rock-5a.dtb")},
			DtoPaths: []FileAsset{
				{Path: filepath.Join(fmt.Sprintf(constants.DtbsAssetPath, "arm64"), "rockchip", "overlay", "rk3588-uart7-m2.dtbo")},
			},
		},
		Output: Output{
			Kind:      OutKindImage,
			OutFormat: OutFormatXZ,
			ImageOptions: &ImageOptions{
				DiskSize: DefaultRAWDiskSize,
				DiskFormat: DiskFormatRaw,
			},
		},
	},
	constants.BoardRock5b: {
		Arch:       "arm64",
		Platform:   constants.PlatformMetal,
		Board:      constants.BoardRock5b,
		SecureBoot: pointer.To(false),
		Input: Input{
			DtbPath: FileAsset{Path: filepath.Join(fmt.Sprintf(constants.DtbsAssetPath, "arm64"), "rockchip", "rk3588-rock-5b.dtb")},
			DtoPaths: []FileAsset{
				{Path: filepath.Join(fmt.Sprintf(constants.DtbsAssetPath, "arm64"), "rockchip", "overlay", "rk3588-uart7-m2.dtbo")},
			},
		},
		Output: Output{
			Kind:      OutKindImage,
			OutFormat: OutFormatXZ,
			ImageOptions: &ImageOptions{
				DiskSize:   DefaultRAWDiskSize,
				DiskFormat: DiskFormatRaw,
			},
		},
	},
}
