package metrics

import (
	"strconv"
	"sync"
	"time"

	smi "lyndeviceplugin/lynsmi-interface"
	podresources "lyndeviceplugin/lynxi-exporter/pod_resources"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type gaugeVec struct {
	m      *prometheus.GaugeVec
	labels labels
}

func newGaugeVec(opts prometheus.GaugeOpts, labels labels) gaugeVec {
	return gaugeVec{
		m:      promauto.NewGaugeVec(opts, labels),
		labels: labels,
	}
}

func (m gaugeVec) set(device Props, v float64) {
	m.m.WithLabelValues(m.labels.getValues(device)...).Set(v)
}

func (m gaugeVec) reset() {
	m.m.Reset()
}

var _ Recorder = &DeviceRecorder{}

// Device 定义和记录所有device相关的Prometheus指标
type DeviceRecorder struct {
	lynxiDeviceStates      gaugeVec
	lynxiDeviceMemUsed     gaugeVec
	lynxiDeviceApuUsage    gaugeVec
	lynxiDeviceArmUsage    gaugeVec
	lynxiDeviceVicUsage    gaugeVec
	lynxiDeviceIpeUsage    gaugeVec
	lynxiDeviceCurrentTemp gaugeVec
	lynxiBoardPower        gaugeVec
	smi                    smi.SMI
	podRes                 *podresources.PodResources
}

// NewDeviceRecorder 构造一个DeviceRecorder并初始化指标
func NewDeviceRecorder(smi smi.SMI, podRes *podresources.PodResources) *DeviceRecorder {
	ret := &DeviceRecorder{
		lynxiDeviceStates: newGaugeVec(prometheus.GaugeOpts{
			Name: "lynxi_device_state",
			Help: "The state of the device",
		}, deviceMetricLabels),
		lynxiDeviceMemUsed: newGaugeVec(prometheus.GaugeOpts{
			Name: "lynxi_device_mem_used",
			Help: "The memory used of the device with unit KB",
		}, memLabels),
		lynxiDeviceApuUsage: newGaugeVec(prometheus.GaugeOpts{
			Name: "lynxi_device_apu_usage",
			Help: "The apu usage of the device with unit %",
		}, deviceMetricLabels),
		lynxiDeviceArmUsage: newGaugeVec(prometheus.GaugeOpts{
			Name: "lynxi_device_arm_usage",
			Help: "The arm usage of the device with unit %",
		}, deviceMetricLabels),
		lynxiDeviceVicUsage: newGaugeVec(prometheus.GaugeOpts{
			Name: "lynxi_device_vic_usage",
			Help: "The vic usage of the device with unit %",
		}, deviceMetricLabels),
		lynxiDeviceIpeUsage: newGaugeVec(prometheus.GaugeOpts{
			Name: "lynxi_device_ipe_usage",
			Help: "The ipe usage of the device with unit FPS",
		}, deviceMetricLabels),
		lynxiDeviceCurrentTemp: newGaugeVec(prometheus.GaugeOpts{
			Name: "lynxi_device_current_temp",
			Help: "The current temperature of the device with unit ℃",
		}, deviceMetricLabels),
		lynxiBoardPower: newGaugeVec(prometheus.GaugeOpts{
			Name: "lynxi_board_power",
			Help: "The power of the board with unit mW. HM100 is unsupported, and the value is always 0",
		}, boardLabels),
		smi:    smi,
		podRes: podRes,
	}
	return ret
}

const (
	labelProductName  = "ProductName"
	labelManufacturer = "Manufacturer"
	labelMountTime    = "MountTime"
	labelBoardID      = "BoardID"
	labelSerialNumber = "SerialNumber"
	labelModel        = "Model"
	labelID           = "ID"
	labelUUID         = "UUID"
	labelMemTotal     = "MemTotal"
	labelPod          = "owner_pod"
	labelContainer    = "owner_container"
	labelNamespace    = "owner_namespace"
)

type labels []string

var boardLabels labels = labels{labelProductName, labelManufacturer, labelMountTime, labelBoardID, labelSerialNumber}
var deviceMetricLabels labels = append(boardLabels, labels{labelModel, labelID, labelUUID, labelPod, labelNamespace, labelContainer}...)
var memLabels labels = append(deviceMetricLabels, labels{labelMemTotal}...)

var mountTime = time.Now().Format(time.RFC3339)

func (ls labels) getValues(device Props) []string {
	ret := make([]string, len(ls))
	for i := range ls {
		ret[i] = ls.getValue(device, i)
	}
	return ret
}

func (ls labels) getValue(device Props, i int) string {
	switch ls[i] {
	case labelProductName:
		return device.Board.ProductName
	case labelManufacturer:
		return device.Board.Brand
	case labelMountTime:
		return mountTime
	case labelBoardID:
		return strconv.Itoa(int(device.Board.ID))
	case labelSerialNumber:
		return device.Board.SerialNumber
	case labelModel:
		return device.Device.Name
	case labelID:
		return device.ID
	case labelUUID:
		return device.Device.UUID
	case labelMemTotal:
		return strconv.Itoa(int(device.Device.MemoryTotal))
	case labelPod:
		return device.Pod
	case labelContainer:
		return device.Container
	case labelNamespace:
		return device.Namespace
	default:
		panic("unknown label")
	}
}

func (m *DeviceRecorder) reset() {
	m.lynxiDeviceStates.reset()
	m.lynxiDeviceMemUsed.reset()
	m.lynxiDeviceApuUsage.reset()
	m.lynxiDeviceArmUsage.reset()
	m.lynxiDeviceVicUsage.reset()
	m.lynxiDeviceIpeUsage.reset()
	m.lynxiDeviceCurrentTemp.reset()
	m.lynxiBoardPower.reset()
}

type Props struct {
	podresources.ResourceOwner
	smi.Props
	ID string
}

func (m *DeviceRecorder) getResoureOwners() (map[string]podresources.ResourceOwner, error) {
	pod_res, err := m.podRes.Get()
	if err != nil {
		return nil, err
	}
	ret := make(map[string]podresources.ResourceOwner)
	for _, res := range pod_res {
		for _, id := range res.IDs {
			ret[id] = res.ResourceOwner
		}
	}
	return ret, err
}

func concurrentExec(fns ...func()) {
	wg := sync.WaitGroup{}
	for _, fn := range fns {
		wg.Add(1)
		go func(fn func()) {
			defer wg.Done()
			fn()
		}(fn)
	}
	wg.Wait()
}

func (m *DeviceRecorder) Record() error {
	var devices smi.AllProps
	var id2ResOwner map[string]podresources.ResourceOwner
	concurrentExec(func() {
		var err error
		devices, err = m.smi.GetDevices()
		GlobalRecorder.logIfError(err)
	}, func() {
		var err error
		id2ResOwner, err = m.getResoureOwners()
		GlobalRecorder.logIfError(err)
	})
	m.reset()
	for i, device_ptr := range devices {
		id := strconv.Itoa(i)
		res_owner := id2ResOwner[id]
		if device_ptr != nil {
			device := Props{res_owner, *device_ptr, id}
			m.lynxiDeviceStates.set(device, StateOK)
			m.lynxiDeviceMemUsed.set(device, float64(device.Device.MemoryUsed))
			m.lynxiDeviceApuUsage.set(device, float64(device.Device.ApuUsage))
			m.lynxiDeviceArmUsage.set(device, float64(device.Device.ArmUsage))
			m.lynxiDeviceVicUsage.set(device, float64(device.Device.VicUsage))
			m.lynxiDeviceIpeUsage.set(device, float64(device.Device.IpeUsage))
			m.lynxiDeviceCurrentTemp.set(device, float64(device.Device.Temperature))
			m.lynxiBoardPower.set(device, float64(device.Board.PowerDraw))
		} else {
			m.lynxiDeviceStates.set(Props{res_owner, smi.Props{}, id}, StateErr)
		}
	}
	return nil
}
