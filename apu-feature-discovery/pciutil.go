package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"
)

// LynxiPCI interface allows us to get a list of all LYNXI PCI devices
type LynxiPCI interface {
	Devices() ([]*PCIDevice, error)
}

// PCIDevice represents a single PCI device
type PCIDevice struct {
	Path    string
	Address string
	Class   string
	Vendor  string
}

const (
	// PciDevicesRoot represents base path for all pci devices under sysfs
	PciDevicesRoot = "/sys/bus/pci/devices"
	// PciLynxiVendorID represents PCI vendor id for Nvidia
	PciLynxiVendorID = "0x1e9f"
)

// LynxiPCILib implements the LynxiPCI interface
type LynxiPCILib struct{}

// NewLynxiPCILib returns an instance of LynxiPCILib implementing the LynxiPCI interface
func NewLynxiPCILib() LynxiPCI {
	return &LynxiPCILib{}
}

// Devices returns all PCI devices on the system
func (p *LynxiPCILib) Devices() ([]*PCIDevice, error) {
	deviceDirs, err := ioutil.ReadDir(PciDevicesRoot)
	if err != nil {
		return nil, fmt.Errorf("Unable to read PCI bus devices: %v", err)
	}

	var devices []*PCIDevice
	for _, deviceDir := range deviceDirs {
		devicePath := path.Join(PciDevicesRoot, deviceDir.Name())
		address := deviceDir.Name()

		vendor, err := ioutil.ReadFile(path.Join(devicePath, "vendor"))
		if err != nil {
			return nil, fmt.Errorf("Unable to read PCI device vendor id for %s: %v", address, err)
		}

		if strings.TrimSpace(string(vendor)) != PciLynxiVendorID {
			continue
		}

		class, err := ioutil.ReadFile(path.Join(devicePath, "class"))
		if err != nil {
			return nil, fmt.Errorf("Unable to read PCI device class for %s: %v", address, err)
		}

		device := &PCIDevice{
			Path:    devicePath,
			Address: address,
			Vendor:  strings.TrimSpace(string(vendor)),
			Class:   string(class)[0:4],
		}

		devices = append(devices, device)
	}

	return devices, nil
}

