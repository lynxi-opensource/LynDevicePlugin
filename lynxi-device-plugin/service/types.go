package service

import (
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

type Service interface {
	pluginapi.DevicePluginServer
	GetResourceName() string
	GetOptions() *pluginapi.DevicePluginOptions
}
