// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"context"
	"log"
	"os"

	"github.com/siderolabs/go-blockdevice/blockdevice/probe"

	"github.com/siderolabs/talos/internal/app/machined/pkg/runtime/v1alpha1/bootloader/grub"
	"github.com/siderolabs/talos/internal/pkg/meta"
	"github.com/siderolabs/talos/internal/pkg/mount"
	"github.com/siderolabs/talos/pkg/machinery/constants"
)

func revertBootloader() {
	if err := revertBootloadInternal(); err != nil {
		log.Printf("failed to revert bootloader: %s", err)
	}
}

//nolint:gocyclo
func revertBootloadInternal() error {
	metaState, err := meta.New(context.Background(), nil)
	if err != nil {
		if os.IsNotExist(err) {
			// no META, no way to revert
			return nil
		}

		return err
	}

	label, ok := metaState.ReadTag(meta.Upgrade)
	if !ok {
		return nil
	}

	if label == "" {
		if _, err = metaState.DeleteTag(context.Background(), meta.Upgrade); err != nil {
			return err
		}

		return metaState.Flush()
	}

	log.Printf("reverting failed upgrade, switching to %q", label)

	// attempt to mount BOOT partition directly without using other code paths, as they rely on Runtime
	dev, err := probe.GetDevWithPartitionName(constants.BootPartitionLabel)
	if os.IsNotExist(err) {
		// no BOOT partition???
		return nil
	}

	if err != nil {
		return err
	}

	defer dev.Close() //nolint:errcheck

	mp, err := mount.SystemMountPointForLabel(dev.BlockDevice, constants.BootPartitionLabel)
	if err != nil {
		return err
	}

	if mp == nil {
		return nil
	}

	alreadyMounted, err := mp.IsMounted()
	if err != nil {
		return err
	}

	if !alreadyMounted {
		if err = mp.Mount(); err != nil {
			return err
		}

		defer mp.Unmount() //nolint:errcheck
	}

	conf, err := grub.Read(grub.ConfigPath)
	if err != nil {
		return err
	}

	if conf == nil {
		return nil
	}

	bootEntry, err := grub.ParseBootLabel(label)
	if err != nil {
		return err
	}

	if conf.Default != bootEntry {
		conf.Default, conf.Fallback = bootEntry, conf.Default
		if err = conf.Write(grub.ConfigPath); err != nil {
			return err
		}
	}

	if _, err = metaState.DeleteTag(context.Background(), meta.Upgrade); err != nil {
		return err
	}

	return metaState.Flush()
}
