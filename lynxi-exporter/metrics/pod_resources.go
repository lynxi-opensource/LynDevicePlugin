package metrics

import (
	smi "lyndeviceplugin/lynsmi-service-client-go"
	podresources "lyndeviceplugin/lynxi-exporter/pod_resources"
	"sort"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var _ Recorder = &PodContainerRecorder{}

// States 定义和记录所有状态相关的Prometheus指标
type PodContainerRecorder struct {
	lynxiPodContainerDeviceCount *prometheus.GaugeVec
	lynxiPodContainerHP280Count  *prometheus.GaugeVec
	deviceID2UUID                map[string]string
	smi                          smi.LynSMI
	podRes                       *podresources.PodResources
}

// NewStatesRecorder 构造一个StatesRecorder并初始化指标
func NewPodContainerRecorder(smi smi.LynSMI, podRes *podresources.PodResources) *PodContainerRecorder {
	ret := &PodContainerRecorder{
		lynxiPodContainerDeviceCount: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "lynxi_pod_container_device_count",
			Help: "The device ids and number of devices for each pod container.",
		}, labelsForPodContainerDevice()),
		lynxiPodContainerHP280Count: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "lynxi_pod_container_hp280_count",
			Help: "The device ids and number of devices for each pod container.",
		}, labelsForPodContainerHP280()),
		deviceID2UUID: make(map[string]string),
		smi:           smi,
		podRes:        podRes,
	}
	return ret
}

func labelsForPodContainerDevice() []string {
	return []string{"owner_pod", "owner_container", "owner_namespace", "device_ids", "uuids"}
}

func labelsForPodContainerHP280() []string {
	return []string{"owner_pod", "owner_container", "owner_namespace", "serial_numbers"}
}

func (m PodContainerRecorder) updateUUIDs() {
	deviceInfos, err := m.smi.GetDevices()
	if err != nil {
		GlobalRecorder.LogError(err)
	}
	for i, deviceInfo := range deviceInfos {
		if info, err := deviceInfo.Get(); err == nil && info != nil {
			m.deviceID2UUID[strconv.Itoa(int(i))] = info.Device.UUID
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

type StringNumberSlice []string

func (a StringNumberSlice) Len() int      { return len(a) }
func (a StringNumberSlice) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a StringNumberSlice) Less(i, j int) bool {
	if len(a[i]) == len(a[j]) {
		return a[i] < a[j]
	}
	return len(a[i]) < len(a[j])
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
			if len(res.IDs) != 0 {
				sort.Sort(StringNumberSlice(res.IDs))
				m.lynxiPodContainerDeviceCount.WithLabelValues(
					res.Pod, res.Container, res.Namespace,
					strings.Join(res.IDs, ","),
					strings.Join(m.getUUIDs(res.IDs), ",")).Set(float64(len(res.IDs)))
			}
			if len(res.HP280s) != 0 {
				m.lynxiPodContainerHP280Count.WithLabelValues(
					res.Pod, res.Container, res.Namespace,
					strings.Join(res.HP280s, ",")).Set(float64(len(res.HP280s)))
			}
		}
	}
	return nil
}
