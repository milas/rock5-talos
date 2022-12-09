// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package makefs

import (
	"fmt"

	"github.com/siderolabs/go-cmd/pkg/cmd"
)

const (
	// FilesystemTypeExt4 is the filesystem type for ext4.
	FilesystemTypeExt4 = "ext4"
)

// Ext4Grow expands an ext4 filesystem to the maximum possible. The partition
// MUST be mounted, or this will fail.
func Ext4Grow(partname string) error {
	_, err := cmd.Run("resize2fs", partname)

	return err
}

// Ext4Repair repairs an ext4 filesystem on the specified partition.
func Ext4Repair(partname, fsType string) error {
	if fsType != FilesystemTypeExt4 {
		return fmt.Errorf("unsupported filesystem type: %s", fsType)
	}

	_, err := cmd.Run("fsck", partname)

	return err
}

// Ext4 creates a ext4 filesystem on the specified partition.
func Ext4(partname string, setters ...Option) error {
	if partname == "" {
		return fmt.Errorf("missing path to disk")
	}

	opts := NewDefaultOptions(setters...)

	var args []string
	if opts.Force {
		args = append(args, "-F")
	}
	if opts.Label != "" {
		args = append(args, "-L", opts.Label)
	}
	args = append(args, partname)

	_, err := cmd.Run("mkfs.ext4", args...)

	return err
}
