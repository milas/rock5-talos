package extlinux

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/siderolabs/talos/internal/app/machined/pkg/runtime/v1alpha1/bootloader/options"
	"github.com/siderolabs/talos/pkg/imager/utils"
	"github.com/siderolabs/talos/pkg/machinery/constants"
)

// BootLabel represents a boot label, e.g. A or B.
type BootLabel string

const (
	// ConfigPath is the path to the extlinux config.
	ConfigPath = constants.BootMountPoint + "/" + constants.ExtlinuxAsset
	// BootA is a bootloader label.
	BootA BootLabel = "A"
	// BootB is a bootloader label.
	BootB BootLabel = "B"
)

const confTemplate = `
label Talos
    kernel {{ .Kernel }}
    initrd {{ .Initrd }}
    devicetreedir {{ .DeviceTreeDir }}
    fdt {{ .Fdt }}
    fdtoverlays {{ .FdtOverlays }}
    append rootwait keepinitrd retain_initrd {{ .Cmdline }}
`

type Config struct {
	Default  BootLabel
	Fallback BootLabel
	Entry    *Entry
}

type Entry struct {
	Kernel        string
	Cmdline       string
	Initrd        string
	DeviceTreeDir string
	Fdt           string
	FdtOverlays   string
}

func NewConfig() *Config {
	return &Config{Default: BootA}
}

func (c *Config) Install(options options.InstallOptions) error {
	if err := c.flip(); err != nil {
		return err
	}

	options.BootAssets.FillDefaults(options.Arch)

	instructions := []utils.CopyInstruction{
		utils.SourceDestination(options.BootAssets.KernelPath, filepath.Join(constants.BootMountPoint, string(c.Default), constants.KernelAsset)),
		utils.SourceDestination(options.BootAssets.InitramfsPath, filepath.Join(constants.BootMountPoint, string(c.Default), constants.InitramfsAsset)),
	}

	if options.BootAssets.DtbPath == "" {
		return errors.New("no device tree .dtb specified")
	}

	dtBasePath := filepath.Join("/", string(c.Default), "dtbs")
	dtbPath := filepath.Join(dtBasePath, filepath.Base(options.BootAssets.DtbPath))
	instructions = append(instructions, utils.SourceDestination(
		options.BootAssets.DtbPath,
		filepath.Join(constants.BootMountPoint, dtbPath),
	))

	dtoPaths := make([]string, len(options.BootAssets.DtoPaths))
	for i := range options.BootAssets.DtoPaths {
		dtoPaths[i] = filepath.Join(dtBasePath, "overlay", filepath.Base(options.BootAssets.DtoPaths[i]))
		instructions = append(instructions, utils.SourceDestination(
			options.BootAssets.DtoPaths[i],
			filepath.Join(constants.BootMountPoint, dtoPaths[i])),
		)
	}

	if err := utils.CopyFiles(options.Printf, instructions...); err != nil {
		return err
	}

	c.Entry = &Entry{
		Kernel:        filepath.Join("/", string(c.Default), constants.KernelAsset),
		Initrd:        filepath.Join("/", string(c.Default), constants.InitramfsAsset),
		Cmdline:       options.Cmdline,
		DeviceTreeDir: dtBasePath,
		Fdt:           dtbPath,
		FdtOverlays:   strings.Join(dtoPaths, " "),
	}

	if err := c.Write(ConfigPath, options.Printf); err != nil {
		return err
	}

	return nil
}

func (c *Config) Revert(_ context.Context) error {
	//TODO implement me
	return nil
}

func (c *Config) PreviousLabel() string {
	return string(c.Fallback)
}

func (c *Config) UEFIBoot() bool {
	return false
}

func (c *Config) Write(path string, printf func(string, ...any)) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModeDir); err != nil {
		return err
	}

	wr := new(bytes.Buffer)

	err := c.Encode(wr)
	if err != nil {
		return err
	}

	printf("writing %s to disk", path)

	return os.WriteFile(path, wr.Bytes(), 0o600)
}

func (c *Config) Encode(wr io.Writer) error {
	t := template.Must(template.New("extlinux").Parse(confTemplate))
	return t.Execute(wr, c.Entry)
}

// flipBootLabel flips the boot label.
func flipBootLabel(e BootLabel) (BootLabel, error) {
	switch e {
	case BootA:
		return BootB, nil
	case BootB:
		return BootA, nil
	//case BootReset:
	//	fallthrough
	default:
		return "", fmt.Errorf("invalid entry: %s", e)
	}
}

// Flip flips the default boot label.
func (c *Config) flip() error {
	if c.Entry == nil {
		return nil
	}

	current := c.Default

	next, err := flipBootLabel(c.Default)
	if err != nil {
		return err
	}

	c.Default = next
	c.Fallback = current

	return nil
}
