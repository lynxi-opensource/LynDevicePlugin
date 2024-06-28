// package service 实现了与kubelet通信所需的grpc接口
package service

import (
	"context"
	"fmt"
	"log"
	"lyndeviceplugin/lynsmi-service-client-go"
	"lyndeviceplugin/lynxi-device-plugin/allocator"
	"math"
	"sort"
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

func getBestDeviceList(scoreMap map[[2]int][2]int32, must []int, availible []int, size int32) []int {
	if size == 0 {
		return must
	}
	if len(availible) == int(size) {
		return availible
	}

	ret := make([]int, 0, len(must))
	ret = append(ret, must...)
	maxScore := int32(-1)
	minDist := int32(math.MaxInt32)
	maxScoreDeviceList := make([]int, 0)
	count := 0
	forEachDeviceList(scoreMap, 0, 0, ret, availible, size, func(score int32, dist int32, selected []int) bool {
		if score > maxScore || (score == maxScore && dist < minDist) {
			maxScore = score
			minDist = dist
			maxScoreDeviceList = append(maxScoreDeviceList[:0], selected...)
		}
		count++
		if count%1000000 == 0 {
			log.Println("count ", count)
		}
		return true
	})
	log.Println("selected", maxScoreDeviceList, "score", maxScore, "dist", minDist)
	return maxScoreDeviceList
}

func getBestDevice(scoreMap map[[2]int][2]int32, selected []int, availible []int) (int, int32, int32) {
	maxScore := int32(-1)
	minDist := int32(math.MaxInt32)
	bestDeviceIndex := -1
	for i, availibleDevice := range availible {
		score := int32(0)
		dist := int32(0)
		for _, selectedDevice := range selected {
			scoreDist := scoreMap[[2]int{availibleDevice, selectedDevice}]
			score += scoreDist[0]
			dist += scoreDist[1]
		}
		if score > maxScore || score == maxScore && dist < minDist {
			maxScore = score
			minDist = dist
			bestDeviceIndex = i
		}
	}
	return bestDeviceIndex, maxScore, minDist
}

func remove(a []int, index int) []int {
	if index != len(a)-1 {
		copy(a[index:], a[index+1:])
	}
	return a[:len(a)-1]
}

func getDeviceListByGreedyPolicy(scoreMap map[[2]int][2]int32, must []int, availible []int, size int32) []int {
	if size == 0 {
		return must
	}
	if len(availible) == int(size) {
		return availible
	}

	sort.Ints(availible)
	if len(must) != 0 {
		selected := make([]int, len(must))
		copy(selected, must)
		nextAvailable := make([]int, len(availible))
		copy(nextAvailable, availible)
		for i := 0; i < int(size); i++ {
			deviceIndex, _, _ := getBestDevice(scoreMap, selected, nextAvailable)
			device := nextAvailable[deviceIndex]
			nextAvailable = remove(nextAvailable, deviceIndex)
			selected = append(selected, device)
		}
		return selected
	}

	maxScore := int32(-1)
	minDist := int32(math.MaxInt32)
	maxScoreDeviceList := make([]int, 0)
	for i := 0; i < len(availible); i++ {
		selected := make([]int, 1)
		selected[0] = availible[i]
		nextAvailable := make([]int, len(availible))
		copy(nextAvailable, availible)
		nextAvailable = remove(nextAvailable, i)
		score := int32(0)
		dist := int32(0)
		for i := 0; i < int(size-1); i++ {
			deviceIndex, deviceScore, deviceDist := getBestDevice(scoreMap, selected, nextAvailable)
			device := nextAvailable[deviceIndex]
			score += deviceScore
			dist += deviceDist
			nextAvailable = remove(nextAvailable, deviceIndex)
			selected = append(selected, device)
		}
		if score > maxScore || score == maxScore && dist < minDist {
			maxScore = score
			minDist = dist
			maxScoreDeviceList = selected
		}
	}
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

func getScoreMap(smi lynsmi.LynSMI) (map[[2]int][2]int32, error) {
	list, err := smi.GetDeviceTopologyList()
	if err != nil {
		return nil, fmt.Errorf("smi GetDeviceTopologyList err: %w", err)
	}
	if list == nil {
		return nil, fmt.Errorf("smi GetDeviceTopologyList return nil")
	}
	scoreMap := make(map[[2]int][2]int32, 0)
	for _, v := range *list {
		attr, err := v.Get()
		if err != nil {
			log.Println("smi GetDeviceTopologyList contain err: ", v)
			continue
		}
		linkScore, dist := getScore(attr.Attr)
		pair := [2]int{int(attr.DevicePair[0]), int(attr.DevicePair[1])}
		score := [2]int32{linkScore, dist}
		scoreMap[pair] = score
		scoreMap[[2]int{pair[1], pair[0]}] = score
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

		must, err := sliceString2Int(containerReq.MustIncludeDeviceIDs)
		if err != nil {
			log.Println("sliceString2Int: ", containerReq.MustIncludeDeviceIDs, "; err: ", err)
			continue
		}
		available, err := sliceString2Int(availableDeviceIDs)
		if err != nil {
			log.Println("sliceString2Int: ", availableDeviceIDs, "; err: ", err)
			continue
		}
		maxScoreDeviceList := getDeviceListByGreedyPolicy(scoreMap, must, available, containerReq.AllocationSize)
		maxScoreDeviceListString := sliceInt2String(maxScoreDeviceList)
		containerResponses[i] = &pluginapi.ContainerPreferredAllocationResponse{
			DeviceIDs: maxScoreDeviceListString,
		}
		log.Println("allocation success:", maxScoreDeviceList)
		allocatedDeviceIDs = append(allocatedDeviceIDs, maxScoreDeviceListString...)
	}
	return &pluginapi.PreferredAllocationResponse{ContainerResponses: containerResponses}, nil
}

func sliceString2Int(src []string) (ret []int, err error) {
	ret = make([]int, 0, len(src))
	for _, v := range src {
		n, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		ret = append(ret, n)
	}
	return ret, nil
}

func sliceInt2String(src []int) (ret []string) {
	ret = make([]string, 0, len(src))
	for _, v := range src {
		n := strconv.Itoa(v)
		ret = append(ret, n)
	}
	return ret
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

func forEachDeviceList(scoreMap map[[2]int][2]int32, score int32, dist int32, selected []int, availible []int, size int32, cb func(score int32, dist int32, selected []int) bool) {
	if size == 0 {
		cb(score, dist, selected)
		return
	}
	length := len(selected)
	for i, v := range availible {
		selected = append(selected, v)
		for _, id := range selected {
			v := scoreMap[[2]int{id, v}]
			score += v[0]
			dist += v[1]
		}
		forEachDeviceList(scoreMap, score, dist, selected, availible[i+1:], size-1, cb)
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
