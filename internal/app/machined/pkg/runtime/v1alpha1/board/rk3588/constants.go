package rk3588

import (
	"fmt"
	"path/filepath"

	"github.com/siderolabs/talos/pkg/machinery/constants"
)

var DeviceTreeOverlays = []string{
	filepath.Join(fmt.Sprintf(constants.DtbsAssetPath, "arm64"), "rockchip", "overlay", "rk3588-uart7-m2.dtbo"),
}
