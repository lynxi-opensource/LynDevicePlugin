// package service 实现了与kubelet通信所需的grpc接口
package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"lyndeviceplugin/lynxi-k8s-device-plugin/allocator"
	"lyndeviceplugin/smi"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

var _ DeviceGetter = &smiMock{}

type smiMock struct {
	devices []smi.Device
	mtx     sync.Mutex
}

func (m *smiMock) GetDevices() ([]smi.Device, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	var err error
	for _, d := range m.devices {
		if !d.IsOn {
			err = errors.New("device is off")
		}
	}
	return m.devices, err
}

func (m *smiMock) setDevices(devices []smi.Device) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.devices = m.devices[:0]
	m.devices = append(m.devices, devices...)
}

var _ allocator.Allocator = allocatorMock{}

type allocatorMock struct{}

func (m allocatorMock) Allocate(*pluginapi.ContainerAllocateRequest) *pluginapi.ContainerAllocateResponse {
	return &pluginapi.ContainerAllocateResponse{}
}

var _ pluginapi.DevicePlugin_ListAndWatchServer = &sendMocker{}

type sendMocker struct {
	callback func(resp *pluginapi.ListAndWatchResponse)
}

func (m *sendMocker) Send(resp *pluginapi.ListAndWatchResponse) error {
	m.callback(resp)
	return nil
}
func (m *sendMocker) SetHeader(metadata.MD) error  { panic("unimplemented") }
func (m *sendMocker) SendHeader(metadata.MD) error { panic("unimplemented") }
func (m *sendMocker) SetTrailer(metadata.MD)       { panic("unimplemented") }
func (m *sendMocker) Context() context.Context     { panic("unimplemented") }
func (m *sendMocker) SendMsg(interface{}) error    { panic("unimplemented") }
func (m *sendMocker) RecvMsg(interface{}) error    { panic("unimplemented") }

func isSMIDevicesAndPluginDevicesEqual(smiDevices []smi.Device, pluginDevices []*pluginapi.Device) bool {
	if len(smiDevices) != len(pluginDevices) {
		return false
	}
	isHealthy := func(isOn bool) string {
		if isOn {
			return pluginapi.Healthy
		}
		return pluginapi.Unhealthy
	}
	for _, sd := range smiDevices {
		found := false
		for _, pd := range pluginDevices {
			if pd.ID == strconv.Itoa(sd.ID) {
				if pd.Health != isHealthy(sd.IsOn) {
					return false
				}
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func TestService_ListAndWatch(t *testing.T) {
	interval := time.Millisecond * 100
	smm := &smiMock{}
	var devices = []smi.Device{
		{DeviceInfo: smi.DeviceInfo{ID: 0, IsOn: true}},
		{DeviceInfo: smi.DeviceInfo{ID: 1, IsOn: true}},
		{DeviceInfo: smi.DeviceInfo{ID: 2, IsOn: true}},
	}
	smm.setDevices(devices)
	svc := NewService(smm, allocatorMock{}, make(chan<- error), interval)
	respChan := make(chan *pluginapi.ListAndWatchResponse)
	go func() {
		assert.Nil(t, svc.ListAndWatch(nil, &sendMocker{callback: func(resp *pluginapi.ListAndWatchResponse) {
			respChan <- resp
		}}))
	}()
	// 立即返回devices
	smiDevices, _ := smm.GetDevices()
	pluginDevices := (<-respChan).Devices
	assert.True(t, isSMIDevicesAndPluginDevicesEqual(smiDevices, pluginDevices), fmt.Sprintf("expect: %v, actual: %v", smiDevices, pluginDevices))
	// 改变device状态
	devices[1].IsOn = false
	smm.setDevices(devices)
	smiDevices, _ = smm.GetDevices()
	pluginDevices = (<-respChan).Devices
	assert.True(t, isSMIDevicesAndPluginDevicesEqual(smiDevices, pluginDevices), fmt.Sprintf("expect: %v, actual: %v", smiDevices, pluginDevices))
}

func TestService_Allocate(t *testing.T) {
	svc := NewService(&smiMock{}, allocatorMock{}, make(chan<- error), time.Millisecond*100)
	deviceIDs := []string{"0", "1"}
	reqs := []*pluginapi.ContainerAllocateRequest{
		{DevicesIDs: deviceIDs},
		{DevicesIDs: deviceIDs},
		{DevicesIDs: deviceIDs},
	}
	resps, err := svc.Allocate(context.Background(), &pluginapi.AllocateRequest{ContainerRequests: reqs})
	assert.Nil(t, err)
	assert.Equal(t, len(reqs), len(resps.ContainerResponses))
}
