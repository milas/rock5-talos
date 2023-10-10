// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package profile

import (
	"github.com/siderolabs/go-pointer"
	"github.com/siderolabs/talos/pkg/machinery/constants"
)

const (
	mib = 1024 * 1024

	// MinRAWDiskSize is the minimum size disk we can create. Used for metal images.
	MinRAWDiskSize = 8192 * mib

	// DefaultRAWDiskSize is the value we use for any non-metal images by default.
	DefaultRAWDiskSize = 8192 * mib
)

// Default describes built-in profiles.
var Default = map[string]Profile{
	// SBCs
	constants.BoardRock5a: {
		Arch:       "arm64",
		Platform:   constants.PlatformMetal,
		Board:      constants.BoardRock5a,
		SecureBoot: pointer.To(false),
		Output: Output{
			Kind:      OutKindImage,
			OutFormat: OutFormatXZ,
			ImageOptions: &ImageOptions{
				DiskSize:   MinRAWDiskSize,
				DiskFormat: DiskFormatRaw,
			},
		},
	},
	constants.BoardRock5b: {
		Arch:       "arm64",
		Platform:   constants.PlatformMetal,
		Board:      constants.BoardRock5b,
		SecureBoot: pointer.To(false),
		Output: Output{
			Kind:      OutKindImage,
			OutFormat: OutFormatXZ,
			ImageOptions: &ImageOptions{
				DiskSize: MinRAWDiskSize,
				DiskFormat: DiskFormatRaw,
			},
		},
	},
}
