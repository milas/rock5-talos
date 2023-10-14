// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package mount

import (
	"fmt"
	"github.com/freddierice/go-losetup/v2"
	"github.com/siderolabs/talos/pkg/machinery/constants"
)

// SquashfsMountPoints returns the mountpoints required to boot the system.
func SquashfsMountPoints(prefix string) (mountpoints *Points, err error) {
	var dev losetup.Device

	dev, err = losetup.Attach("/"+constants.RootfsAsset, 0, true)
	if err != nil {
		return nil, fmt.Errorf("squashfs: attach: %w", err)
	}

	squashfs := NewMountPoints()
	// flags make rock5b unhappy?
	squashfs.Set("squashfs", NewMountPoint(dev.Path(), "/", "squashfs", 0x0, "", WithPrefix(prefix), WithFlags(Shared)))

	return squashfs, nil
}
