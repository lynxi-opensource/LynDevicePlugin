// Package allocator 设置lynxi-docker使用device时所需的参数
package allocator

import (
	"strings"

	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

// Allocator 定义分配设备的接口
type Allocator interface {
	Allocate(*pluginapi.ContainerAllocateRequest) *pluginapi.ContainerAllocateResponse
}

type Alloc struct {
}

func (alloc Alloc) Allocate(req *pluginapi.ContainerAllocateRequest) *pluginapi.ContainerAllocateResponse {
	resp := pluginapi.ContainerAllocateResponse{Envs: make(map[string]string)}
	resp.Envs["LYNXI_VISIBLE_DEVICES"] += strings.Join(req.GetDevicesIDs(), ",")
	return &resp
}
