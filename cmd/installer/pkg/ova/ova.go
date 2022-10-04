// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package ova implements OVA creation.
package ova

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/talos-systems/go-cmd/pkg/cmd"

	"github.com/talos-systems/talos/cmd/installer/pkg"
	"github.com/talos-systems/talos/cmd/installer/pkg/qemuimg"
)

const mfTpl = `SHA256({{ .VMDK }})= {{ .VMDKSHA }}
SHA256({{ .OVF }})= {{ .OVFSHA }}
`

// OVF format reference: https://www.dmtf.org/standards/ovf.
//
//nolint:lll
const ovfTpl = `<?xml version="1.0" encoding="UTF-8"?>
<!--Generated by VMware ovftool 4.3.0 (build-7948156), UTC time: 2019-10-31T01:41:10.540841Z-->
<!-- Edited by Talos -->
<Envelope vmw:buildId="build-7948156" xmlns="http://schemas.dmtf.org/ovf/envelope/1" xmlns:cim="http://schemas.dmtf.org/wbem/wscim/1/common" xmlns:ovf="http://schemas.dmtf.org/ovf/envelope/1" xmlns:rasd="http://schemas.dmtf.org/wbem/wscim/1/cim-schema/2/CIM_ResourceAllocationSettingData" xmlns:vmw="http://www.vmware.com/schema/ovf" xmlns:vssd="http://schemas.dmtf.org/wbem/wscim/1/cim-schema/2/CIM_VirtualSystemSettingData" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <References>
    <File ovf:href="{{ .VMDK }}" ovf:id="file1" ovf:size="{{ .Size }}"/>
  </References>
  <DiskSection>
    <Info>Virtual disk information</Info>
    <Disk ovf:capacity="{{ .Capacity }}" ovf:capacityAllocationUnits="byte * 2^20" ovf:diskId="vmdisk1" ovf:fileRef="file1" ovf:format="http://www.vmware.com/interfaces/specifications/vmdk.html#streamOptimized"/>
  </DiskSection>
  <NetworkSection>
    <Info>The list of logical networks</Info>
    <Network ovf:name="VM Network">
      <Description>The VM Network network</Description>
    </Network>
  </NetworkSection>
  <VirtualSystem ovf:id="vm">
    <Info>A virtual machine</Info>
    <Name>talos</Name>
    <ProductSection ovf:required="false">
      <Info>Talos Virtual Appliance</Info>
      <Property ovf:userConfigurable="true" ovf:type="string"
                ovf:key="talos.config" ovf:value="">
        <Label>Talos config data</Label>
        <Description>Inline Talos config</Description>
      </Property>
    </ProductSection>
    <OperatingSystemSection ovf:id="101" vmw:osType="other3xLinux64Guest">
      <Info>The kind of installed guest operating system</Info>
    </OperatingSystemSection>
    <VirtualHardwareSection ovf:transport="com.vmware.guestInfo">
      <Info>Virtual hardware requirements</Info>
      <System>
        <vssd:ElementName>Virtual Hardware Family</vssd:ElementName>
        <vssd:InstanceID>0</vssd:InstanceID>
        <vssd:VirtualSystemIdentifier>talos</vssd:VirtualSystemIdentifier>
        <vssd:VirtualSystemType>vmx-13</vssd:VirtualSystemType>
      </System>
      <Item>
        <rasd:AllocationUnits>hertz * 10^6</rasd:AllocationUnits>
        <rasd:Description>Number of Virtual CPUs</rasd:Description>
        <rasd:ElementName>2 virtual CPU(s)</rasd:ElementName>
        <rasd:InstanceID>1</rasd:InstanceID>
        <rasd:ResourceType>3</rasd:ResourceType>
        <rasd:VirtualQuantity>2</rasd:VirtualQuantity>
      </Item>
      <Item>
        <rasd:AllocationUnits>byte * 2^20</rasd:AllocationUnits>
        <rasd:Description>Memory Size</rasd:Description>
        <rasd:ElementName>2048MB of memory</rasd:ElementName>
        <rasd:InstanceID>2</rasd:InstanceID>
        <rasd:ResourceType>4</rasd:ResourceType>
        <rasd:VirtualQuantity>2048</rasd:VirtualQuantity>
      </Item>
      <Item>
        <rasd:Address>0</rasd:Address>
        <rasd:Description>SCSI Controller</rasd:Description>
        <rasd:ElementName>scsiController0</rasd:ElementName>
        <rasd:InstanceID>3</rasd:InstanceID>
        <rasd:ResourceSubType>VirtualSCSI</rasd:ResourceSubType>
        <rasd:ResourceType>6</rasd:ResourceType>
      </Item>
      <Item>
        <rasd:AddressOnParent>0</rasd:AddressOnParent>
        <rasd:ElementName>disk0</rasd:ElementName>
        <rasd:HostResource>ovf:/disk/vmdisk1</rasd:HostResource>
        <rasd:InstanceID>4</rasd:InstanceID>
        <rasd:Parent>3</rasd:Parent>
        <rasd:ResourceType>17</rasd:ResourceType>
      </Item>
      <Item>
        <rasd:AddressOnParent>2</rasd:AddressOnParent>
        <rasd:AutomaticAllocation>true</rasd:AutomaticAllocation>
        <rasd:Connection>VM Network</rasd:Connection>
        <rasd:Description>VmxNet3 ethernet adapter on &quot;VM Network&quot;</rasd:Description>
        <rasd:ElementName>ethernet0</rasd:ElementName>
        <rasd:InstanceID>5</rasd:InstanceID>
        <rasd:ResourceSubType>VmxNet3</rasd:ResourceSubType>
        <rasd:ResourceType>10</rasd:ResourceType>
        <vmw:Config ovf:required="false" vmw:key="slotInfo.pciSlotNumber" vmw:value="32"/>
        <vmw:Config ovf:required="false" vmw:key="wakeOnLanEnabled" vmw:value="false"/>
        <vmw:Config ovf:required="false" vmw:key="connectable.allowGuestControl" vmw:value="false"/>
      </Item>
      <Item ovf:required="false">
        <rasd:AutomaticAllocation>false</rasd:AutomaticAllocation>
        <rasd:ElementName>video</rasd:ElementName>
        <rasd:InstanceID>6</rasd:InstanceID>
        <rasd:ResourceType>24</rasd:ResourceType>
      </Item>
      <Item ovf:required="false">
        <rasd:AutomaticAllocation>false</rasd:AutomaticAllocation>
        <rasd:ElementName>vmci</rasd:ElementName>
        <rasd:InstanceID>7</rasd:InstanceID>
        <rasd:ResourceSubType>vmware.vmci</rasd:ResourceSubType>
        <rasd:ResourceType>1</rasd:ResourceType>
      </Item>
      <vmw:Config ovf:required="false" vmw:key="tools.syncTimeWithHost" vmw:value="false"/>
      <vmw:Config ovf:required="false" vmw:key="tools.afterPowerOn" vmw:value="true"/>
      <vmw:Config ovf:required="false" vmw:key="tools.afterResume" vmw:value="true"/>
      <vmw:Config ovf:required="false" vmw:key="tools.beforeGuestShutdown" vmw:value="true"/>
      <vmw:Config ovf:required="false" vmw:key="tools.beforeGuestStandby" vmw:value="true"/>
      <vmw:Config ovf:required="false" vmw:key="tools.toolsUpgradePolicy" vmw:value="manual"/>
      <vmw:Config ovf:required="false" vmw:key="powerOpInfo.suspendType" vmw:value="soft"/>
      <vmw:ExtraConfig ovf:required="false" vmw:key="nvram" vmw:value="talos.nvram"/>
    </VirtualHardwareSection>
  </VirtualSystem>
</Envelope>
`

// CreateOVAFromRAW creates an OVA from a RAW disk.
//
//nolint:gocyclo
func CreateOVAFromRAW(name, src, out, arch string) (err error) {
	dir, err := os.MkdirTemp("/tmp", "talos")
	if err != nil {
		return err
	}

	dest := filepath.Join(dir, name+".vmdk")

	if err = qemuimg.Convert("raw", "vmdk", "compat6,subformat=streamOptimized,adapter_type=lsilogic", src, dest); err != nil {
		return err
	}

	f, err := os.Stat(dest)
	if err != nil {
		return err
	}

	size := f.Size()

	ovf, err := renderOVF(name, size, pkg.RAWDiskSize)
	if err != nil {
		return err
	}

	input, err := os.Open(dest)
	if err != nil {
		return err
	}

	vmdkSHA25Sum, err := sha256sum(input)
	if err != nil {
		return err
	}

	ovfSHA25Sum, err := sha256sum(strings.NewReader(ovf))
	if err != nil {
		return err
	}

	mf, err := renderMF(name, vmdkSHA25Sum, ovfSHA25Sum)
	if err != nil {
		return err
	}

	//nolint:errcheck
	defer os.RemoveAll(dir)

	if err = os.WriteFile(filepath.Join(dir, name+".mf"), []byte(mf), 0o666); err != nil {
		return err
	}

	if err = os.WriteFile(filepath.Join(dir, name+".ovf"), []byte(ovf), 0o666); err != nil {
		return err
	}

	if _, err = cmd.Run("tar", "-cvf", filepath.Join(out, fmt.Sprintf("vmware-%s.ova", arch)), "-C", dir, name+".ovf", name+".mf", name+".vmdk"); err != nil {
		return err
	}

	return nil
}

func sha256sum(input io.Reader) (string, error) {
	hash := sha256.New()

	if _, err := io.Copy(hash, input); err != nil {
		return "", err
	}

	sum := hash.Sum(nil)

	return fmt.Sprintf("%x", sum), nil
}

func renderMF(name, vmdkSHA25Sum, ovfSHA25Sum string) (string, error) {
	cfg := struct {
		VMDK    string
		VMDKSHA string
		OVF     string
		OVFSHA  string
	}{
		VMDK:    name + ".vmdk",
		VMDKSHA: vmdkSHA25Sum,
		OVF:     name + ".ovf",
		OVFSHA:  ovfSHA25Sum,
	}

	templ := template.Must(template.New("mf").Parse(mfTpl))

	var buf bytes.Buffer

	if err := templ.Execute(&buf, cfg); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func renderOVF(name string, size, capacity int64) (string, error) {
	cfg := struct {
		VMDK     string
		Size     int64
		Capacity int64
	}{
		VMDK:     name + ".vmdk",
		Size:     size,
		Capacity: capacity,
	}

	templ := template.Must(template.New("ovf").Parse(ovfTpl))

	var buf bytes.Buffer

	if err := templ.Execute(&buf, cfg); err != nil {
		return "", err
	}

	return buf.String(), nil
}
