package rk3588

import (
	"fmt"
	"github.com/siderolabs/talos/pkg/machinery/constants"
	"path/filepath"
)

var DeviceTreeOverlays = []string{
	filepath.Join(fmt.Sprintf(constants.DtbsAssetPath, "arm64"), "rockchip", "overlay", "rk3588-uart7-m2.dtbo"),
}
