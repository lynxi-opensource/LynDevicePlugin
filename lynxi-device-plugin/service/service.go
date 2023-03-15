// package service 实现了与kubelet通信所需的grpc接口
package service

import (
	"context"
	"errors"
	"log"
	"strconv"
	"time"

	"lyndeviceplugin/lynxi-device-plugin/allocator"
	"lyndeviceplugin/smi"

	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

var _ pluginapi.DevicePluginServer = &Service{}

// DeviceGetter 描述获取设备状态的接口，smi.SMI的子集
type DeviceGetter interface {
	// GetDevices 获取所有设备
	GetDevices() ([]smi.Device, error)
}

// Service 实现了pluginapi.DevicePluginServer，提供grpc接口实现
type Service struct {
	allocator allocator.Allocator
	smi       DeviceGetter
	deviceMap map[int]bool
	interval  time.Duration
}

// NewService 构造一个Service
func NewService(smi DeviceGetter, allocator allocator.Allocator, pollInterval time.Duration) *Service {
	return &Service{
		smi:       smi,
		allocator: allocator,
		deviceMap: make(map[int]bool),
		interval:  pollInterval,
	}
}

// GetDevicePluginOptions 返回PreStartRequired和GetPreferredAllocationAvailable选项为false
func (m *Service) GetDevicePluginOptions(context.Context, *pluginapi.Empty) (*pluginapi.DevicePluginOptions, error) {
	return &pluginapi.DevicePluginOptions{}, nil
}

func isHealthy(isOn bool) string {
	if isOn {
		return pluginapi.Healthy
	}
	return pluginapi.Unhealthy
}

func smiDevicesToPluginDevices(devices []smi.Device) (ret []*pluginapi.Device) {
	for _, d := range devices {
		ret = append(ret, &pluginapi.Device{
			ID:     strconv.Itoa(d.ID),
			Health: isHealthy(d.IsOn),
		})
	}
	return
}

// ListAndWatch 返回所有板卡信息
func (m *Service) ListAndWatch(_ *pluginapi.Empty, sender pluginapi.DevicePlugin_ListAndWatchServer) error {
	ticker := time.NewTicker(m.interval)
	log.Println("start send device status")
	for {
		devices, err := m.smi.GetDevices()
		if err != nil {
			log.Println("smi GetDevices err: ", err)
		}
		if err = sender.Send(&pluginapi.ListAndWatchResponse{Devices: smiDevicesToPluginDevices(devices)}); err != nil {
			log.Fatalln("ListAndWatch: send to kubelet err:", err)
		}
		<-ticker.C
	}
}

// GetPreferredAllocation 未实现
func (m *Service) GetPreferredAllocation(context.Context, *pluginapi.PreferredAllocationRequest) (*pluginapi.PreferredAllocationResponse, error) {
	return nil, errors.New("unimplemented")
}

// Allocate 根据设备id分配板卡
func (m Service) Allocate(_ context.Context, req *pluginapi.AllocateRequest) (*pluginapi.AllocateResponse, error) {
	ret := pluginapi.AllocateResponse{}
	for _, r := range req.GetContainerRequests() {
		ret.ContainerResponses = append(ret.ContainerResponses, m.allocator.Allocate(r))
	}
	return &ret, nil
}

// PreStartContainer 没有任何操作
func (m *Service) PreStartContainer(context.Context, *pluginapi.PreStartContainerRequest) (*pluginapi.PreStartContainerResponse, error) {
	return nil, nil
}
