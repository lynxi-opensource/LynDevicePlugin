// package service 实现了与kubelet通信所需的grpc接口
package service

import (
	"context"
	"log"
	"lyndeviceplugin/lynsmi-service-client-go"
	"lyndeviceplugin/lynxi-device-plugin/allocator"
	"strings"
	"time"

	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

const HP280Name = "HP280"
const HP280DeviceCount = 8

// const HP280Name = "HP300"
// const HP280DeviceCount = 3

var _ pluginapi.DevicePluginServer = &HP280AllocationService{}

// HP280AllocationService 实现了pluginapi.DevicePluginServer，提供grpc接口实现
type HP280AllocationService struct {
	smi       lynsmi.LynSMI
	deviceMap map[int]bool
	interval  time.Duration
}

// NewService 构造一个Service
func NewHP280AllocationService(allocator allocator.Allocator, pollInterval time.Duration) *HP280AllocationService {
	return &HP280AllocationService{
		smi:       lynsmi.LynSMI{},
		deviceMap: make(map[int]bool),
		interval:  pollInterval,
	}
}

func (m *HP280AllocationService) GetResourceName() string {
	return "lynxi.com/" + HP280Name
}

func (m *HP280AllocationService) GetOptions() *pluginapi.DevicePluginOptions {
	return &pluginapi.DevicePluginOptions{}
}

// GetDevicePluginOptions 返回PreStartRequired和GetPreferredAllocationAvailable选项为false
func (m *HP280AllocationService) GetDevicePluginOptions(context.Context, *pluginapi.Empty) (*pluginapi.DevicePluginOptions, error) {
	return m.GetOptions(), nil
}

func getBoardDeviceMap(devices lynsmi.PropsMap) map[string][]int {
	boardDeviceMap := make(map[string][]int)
	for i, d := range devices {
		props, err := d.Get()
		if err == nil && props != nil {
			boardDeviceMap[props.Board.SerialNumber] = append(boardDeviceMap[props.Board.SerialNumber], int(i))
		}
	}
	return boardDeviceMap
}

func hp280BoardsToPluginDevices(devices lynsmi.PropsMap) (ret []*pluginapi.Device) {
	boardDeviceMap := getBoardDeviceMap(devices)
	for k, v := range boardDeviceMap {
		ret = append(ret, &pluginapi.Device{
			ID:     k,
			Health: isHealthy(len(v) == HP280DeviceCount),
		})
	}
	return
}

// ListAndWatch 返回所有板卡信息
func (m *HP280AllocationService) ListAndWatch(_ *pluginapi.Empty, sender pluginapi.DevicePlugin_ListAndWatchServer) error {
	ticker := time.NewTicker(m.interval)
	log.Println("start send device status")
	for {
		devices, err := m.smi.GetDevices()
		if err != nil {
			log.Println("smi GetDevices err: ", err)
		}
		if err = sender.Send(&pluginapi.ListAndWatchResponse{Devices: hp280BoardsToPluginDevices(devices)}); err != nil {
			log.Fatalln("ListAndWatch: send to kubelet err:", err)
		}
		<-ticker.C
	}
}

// GetPreferredAllocation
func (m *HP280AllocationService) GetPreferredAllocation(ctx context.Context, req *pluginapi.PreferredAllocationRequest) (*pluginapi.PreferredAllocationResponse, error) {
	return nil, nil
}

// Allocate 根据设备id分配板卡
func (m HP280AllocationService) Allocate(_ context.Context, req *pluginapi.AllocateRequest) (*pluginapi.AllocateResponse, error) {
	devices, err := m.smi.GetDevices()
	if err != nil {
		log.Println("smi GetDevices err: ", err)
		return nil, err
	}
	boardDeviceMap := getBoardDeviceMap(devices)
	ret := pluginapi.AllocateResponse{}
	for _, req := range req.GetContainerRequests() {
		ids := make([]int, 0)
		for _, sn := range req.GetDevicesIDs() {
			boardDevices := boardDeviceMap[sn]
			ids = append(ids, boardDevices...)
		}
		resp := pluginapi.ContainerAllocateResponse{Envs: make(map[string]string)}
		resp.Envs["LYNXI_VISIBLE_DEVICES"] = strings.Join(sliceInt2String(ids), ",")
		ret.ContainerResponses = append(ret.ContainerResponses, &resp)
	}
	return &ret, nil
}

// PreStartContainer 没有任何操作
func (m *HP280AllocationService) PreStartContainer(context.Context, *pluginapi.PreStartContainerRequest) (*pluginapi.PreStartContainerResponse, error) {
	return nil, nil
}
