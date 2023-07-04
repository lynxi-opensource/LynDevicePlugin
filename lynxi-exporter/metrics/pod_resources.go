package metrics

import (
	smi "lyndeviceplugin/lynsmi-interface"
	podresources "lyndeviceplugin/lynxi-exporter/pod_resources"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var _ Recorder = &PodContainerRecorder{}

// States 定义和记录所有状态相关的Prometheus指标
type PodContainerRecorder struct {
	lynxiPodContainerDeviceCount *prometheus.GaugeVec
	deviceID2UUID                map[string]string
	smi                          smi.SMI
	podRes                       *podresources.PodResources
}

// NewStatesRecorder 构造一个StatesRecorder并初始化指标
func NewPodContainerRecorder(smi smi.SMI, podRes *podresources.PodResources) *PodContainerRecorder {
	ret := &PodContainerRecorder{
		lynxiPodContainerDeviceCount: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "lynxi_pod_container_device_count",
			Help: "The device ids and number of devices for each pod container.",
		}, labelsForPodContainer()),
		deviceID2UUID: make(map[string]string),
		smi:           smi,
		podRes:        podRes,
	}
	return ret
}

func labelsForPodContainer() []string {
	return []string{"owner_pod", "owner_container", "owner_namespace", "device_ids", "uuids"}
}

func (m PodContainerRecorder) updateUUIDs() {
	deviceInfos, err := m.smi.GetDevices()
	if err != nil {
		GlobalRecorder.LogError(err)
	}
	for i, deviceInfo := range deviceInfos {
		if deviceInfo != nil {
			m.deviceID2UUID[strconv.Itoa(i)] = deviceInfo.Device.UUID
		}
	}
}

func (m PodContainerRecorder) getUUIDs(deviceIDs []string) (ret []string) {
	for _, id := range deviceIDs {
		uuid := m.deviceID2UUID[id]
		ret = append(ret, uuid)
	}
	return
}

func (m *PodContainerRecorder) Record() error {
	m.updateUUIDs()
	resp, err := m.podRes.Get()
	if err != nil {
		GlobalRecorder.LogError(err)
		return err
	} else {
		m.lynxiPodContainerDeviceCount.Reset()
		for _, res := range resp {
			m.lynxiPodContainerDeviceCount.WithLabelValues(
				res.Pod, res.Container, res.Namespace,
				strings.Join(res.IDs, ","),
				strings.Join(m.getUUIDs(res.IDs), ",")).Set(float64(len(res.IDs)))
		}
	}
	return nil
}
