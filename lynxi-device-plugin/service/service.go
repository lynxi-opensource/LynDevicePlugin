// package service 实现了与kubelet通信所需的grpc接口
package service

import (
	"context"
	"fmt"
	"log"
	"lyndeviceplugin/lynsmi-service-client-go"
	"lyndeviceplugin/lynxi-device-plugin/allocator"
	"math"
	"strconv"
	"time"

	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

var _ pluginapi.DevicePluginServer = &Service{}

// Service 实现了pluginapi.DevicePluginServer，提供grpc接口实现
type Service struct {
	allocator allocator.Allocator
	smi       lynsmi.LynSMI
	deviceMap map[int]bool
	interval  time.Duration
}

// NewService 构造一个Service
func NewService(allocator allocator.Allocator, pollInterval time.Duration) *Service {
	return &Service{
		smi:       lynsmi.LynSMI{},
		allocator: allocator,
		deviceMap: make(map[int]bool),
		interval:  pollInterval,
	}
}

// GetDevicePluginOptions 返回PreStartRequired和GetPreferredAllocationAvailable选项为false
func (m *Service) GetDevicePluginOptions(context.Context, *pluginapi.Empty) (*pluginapi.DevicePluginOptions, error) {
	return &pluginapi.DevicePluginOptions{GetPreferredAllocationAvailable: true}, nil
}

func isHealthy(isOn bool) string {
	if isOn {
		return pluginapi.Healthy
	}
	return pluginapi.Unhealthy
}

func smiDevicesToPluginDevices(devices lynsmi.PropsMap) (ret []*pluginapi.Device) {
	for i, d := range devices {
		props, err := d.Get()
		ret = append(ret, &pluginapi.Device{
			ID:     strconv.Itoa(int(i)),
			Health: isHealthy(err == nil && props != nil),
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

func getBestDeviceList(scoreMap map[[2]string][2]int32, must []string, availible []string, size int32) []string {
	ret := make([]string, 0, len(must))
	ret = append(ret, must...)
	maxScore := int32(-1)
	minDist := int32(math.MaxInt32)
	maxScoreDeviceList := make([]string, 0)
	forEachDeviceList(ret, availible, size, func(selected []string) bool {
		score := int32(0)
		dist := int32(0)
		for _, d1 := range selected {
			for _, d2 := range selected {
				v := scoreMap[[2]string{d1, d2}]
				score += v[0]
				dist += v[1]
			}
		}
		if score > maxScore || (score == maxScore && dist < minDist) {
			maxScore = score
			minDist = dist
			maxScoreDeviceList = append(maxScoreDeviceList[:0], selected...)
		}
		return true
	})
	return maxScoreDeviceList
}

func getScore(attr lynsmi.P2PAttr) (int32, int32) {
	if attr.Mode == lynsmi.P2PLinkPIX {
		return 70, attr.Dist
	}
	if attr.Mode == lynsmi.P2PLinkPXB {
		return 50, attr.Dist
	}
	if attr.Mode == lynsmi.P2PLinkPHB {
		return 30, attr.Dist
	}
	if attr.Mode == lynsmi.P2PLinkSYS {
		return 10, attr.Dist
	}
	if attr.Mode == lynsmi.NonSupport {
		return 0, 0
	}
	log.Println("invalid p2p attr: ", attr)
	return 0, 0
}

func getScoreMap(smi lynsmi.LynSMI) (map[[2]string][2]int32, error) {
	list, err := smi.GetDeviceTopologyList()
	if err != nil {
		return nil, fmt.Errorf("smi GetDeviceTopologyList err: %w", err)
	}
	if list == nil {
		return nil, fmt.Errorf("smi GetDeviceTopologyList return nil")
	}
	scoreMap := make(map[[2]string][2]int32, 0)
	for _, v := range *list {
		attr, err := v.Get()
		if err != nil {
			log.Println("smi GetDeviceTopologyList contain err: ", v)
			continue
		}
		linkScore, dist := getScore(attr.Attr)
		pair := [2]string{strconv.Itoa(int(attr.DevicePair[0])), strconv.Itoa(int(attr.DevicePair[1]))}
		score := [2]int32{linkScore, dist}
		scoreMap[pair] = score
		scoreMap[[2]string{pair[1], pair[0]}] = score
	}
	return scoreMap, nil
}

// GetPreferredAllocation
func (m *Service) GetPreferredAllocation(ctx context.Context, req *pluginapi.PreferredAllocationRequest) (*pluginapi.PreferredAllocationResponse, error) {
	containerRequests := req.GetContainerRequests()
	if containerRequests == nil {
		return &pluginapi.PreferredAllocationResponse{}, nil
	}

	scoreMap, err := getScoreMap(m.smi)
	if err != nil {
		log.Println(err)
		return &pluginapi.PreferredAllocationResponse{}, nil
	}

	allocatedDeviceIDs := make([]string, 0)
	containerResponses := make([]*pluginapi.ContainerPreferredAllocationResponse, len(containerRequests))
	for i, containerReq := range containerRequests {
		if containerReq == nil {
			containerResponses[i] = nil
			continue
		}

		log.Println(
			"ContainerPreferredAllocation:",
			"must:", containerReq.MustIncludeDeviceIDs,
			"avaliable:", containerReq.AvailableDeviceIDs,
			"size:", containerReq.AllocationSize)

		availableDeviceIDs := sliceSub(containerReq.AvailableDeviceIDs, containerReq.MustIncludeDeviceIDs)
		availableDeviceIDs = sliceSub(availableDeviceIDs, allocatedDeviceIDs)
		if len(availableDeviceIDs) < int(containerReq.AllocationSize) {
			log.Println("allocation failed")
			continue
		}

		maxScoreDeviceList := getBestDeviceList(scoreMap, containerReq.MustIncludeDeviceIDs, availableDeviceIDs, containerReq.AllocationSize)
		containerResponses[i] = &pluginapi.ContainerPreferredAllocationResponse{
			DeviceIDs: maxScoreDeviceList,
		}
		log.Println("allocation success:", maxScoreDeviceList)
		allocatedDeviceIDs = append(allocatedDeviceIDs, maxScoreDeviceList...)
	}
	return &pluginapi.PreferredAllocationResponse{ContainerResponses: containerResponses}, nil
}

func sliceSub[T comparable](a []T, b []T) (ret []T) {
	for _, va := range a {
		existInB := false
		for _, vb := range b {
			if va == vb {
				existInB = true
				break
			}
		}
		if !existInB {
			ret = append(ret, va)
		}
	}
	return
}

func forEachDeviceList(selected []string, availible []string, size int32, cb func(selected []string) bool) {
	if size == 0 {
		cb(selected)
	}
	length := len(selected)
	for i, v := range availible {
		selected = append(selected, v)
		forEachDeviceList(selected, availible[i+1:], size-1, cb)
		selected = selected[:length]
	}
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
