package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	gpuPassThroughConfig = []byte{0xde, 0x10, 0x8a, 0x11, 0x07, 0x04, 0x10, 0x00, 0xa1, 0x00, 0x00, 0x03, 0x00, 0xf8, 0x00, 0x00, 0x00, 0x00, 0x00, 0xec, 0x0c, 0x00, 0x00, 0xe0, 0x00, 0x00, 0x00, 0x00, 0x0c, 0x00, 0x00, 0xea, 0x00, 0x00, 0x00, 0x00, 0x01, 0xc1, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xde, 0x10, 0x14, 0x10, 0x00, 0x00, 0x00, 0xee, 0x60, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, 0x01, 0x00, 0x00, 0xde, 0x10, 0x14, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0xce, 0xd6, 0x23, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x68, 0x03, 0x00, 0x08, 0x00, 0x00, 0x00, 0x05, 0x78, 0x81, 0x00, 0x00, 0x70, 0xe6, 0xfe, 0x00, 0x00, 0x00, 0x00, 0x00, 0x43, 0x00, 0x00, 0x10, 0xb4, 0x02, 0x00, 0xe1, 0x8d, 0x64, 0x00, 0x10, 0x29, 0x00, 0x00, 0x03, 0x3d, 0x45, 0x10, 0x00, 0x00, 0x01, 0x11, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x13, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0e, 0x00, 0x00, 0x00, 0x03, 0x00, 0x3e, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x09, 0x00, 0x14, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	vgpuConfig           = []byte{0xde, 0x10, 0xb8, 0x1e, 0x02, 0x05, 0xff, 0x06, 0xa1, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xfc, 0x0c, 0x00, 0x00, 0xd0, 0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0xfa, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xde, 0x10, 0x0f, 0x13, 0x00, 0x00, 0x00, 0x00, 0xd0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0a, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0xce, 0xd6, 0x23, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, 0x00, 0x81, 0x00, 0x00, 0x00, 0xe0, 0xfe, 0x00, 0x00, 0x00, 0x00, 0x4e, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x09, 0x68, 0x1b, 0x56, 0x46, 0x00, 0x16, 0x34, 0x36, 0x30, 0x2e, 0x31, 0x36, 0x00, 0x00, 0x00, 0x00, 0x72, 0x34, 0x36, 0x30, 0x5f, 0x30, 0x30, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
)

// MockNvidiaPCI represents mock of NvidiaPCI interface
type MockNvidiaPCI struct {
	devices []*PCIDevice
}

// Devices returns PCI devices with mocked data
func (p *MockNvidiaPCI) Devices() ([]*PCIDevice, error) {
	return p.devices, nil
}

// NewMockNvidiaPCI initializes and returns mock PCI interface type
func NewMockNvidiaPCI() NvidiaPCI {
	return &MockNvidiaPCI{
		devices: []*PCIDevice{
			{
				Path:    "",
				Address: "passthrough",
				Vendor:  "0x10de",
				Class:   "300",
				Config:  gpuPassThroughConfig,
			},
			{
				Path:    "",
				Address: "vgpu",
				Vendor:  "0x10de",
				Class:   "300",
				Config:  vgpuConfig,
			},
		},
	}
}

func TestGetVendorSpecificCapability(t *testing.T) {
	devices, _ := NewMockNvidiaPCI().Devices()
	for _, device := range devices {
		// check for vendor id
		require.Equal(t, "0x10de", fmt.Sprintf("0x%x", GetWord(device.Config, 0)), "Nvidia PCI Vendor ID")
		// check for vendor specific capability
		capability, err := device.GetVendorSpecificCapability()
		require.NoError(t, err, "Get vendor specific capability from configuration space")
		require.NotZero(t, len(capability), "Vendor capability record")
		if device.Address == "passthrough" {
			require.Equal(t, 20, len(capability), "Vendor capability length for passthrough device")
		}
		if device.Address == "vgpu" {
			require.Equal(t, 27, len(capability), "Vendor capability length for vgpu device")
		}
	}
}
