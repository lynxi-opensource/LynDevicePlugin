package metrics

import (
	smi "lyndeviceplugin/lynsmi-service-client-go"
	podresources "lyndeviceplugin/lynxi-exporter/pod_resources"
	"strconv"
	"sync"
	"time"

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

func (m gaugeVec) setNoProps(device Props, v float64) {
	m.m.WithLabelValues(m.labels.getNoProps(device)...).Set(v)
}

func (m gaugeVec) reset() {
	m.m.Reset()
}

var _ Recorder = &DeviceRecorder{}

// Device 定义和记录所有device相关的Prometheus指标
type DeviceRecorder struct {
	lynxiDeviceState              gaugeVec
	lynxiDeviceMemUsed            gaugeVec
	lynxiDeviceApuUsage           gaugeVec
	lynxiDeviceArmUsage           gaugeVec
	lynxiDeviceVicUsage           gaugeVec
	lynxiDeviceIpeUsage           gaugeVec
	lynxiDeviceCurrentTemp        gaugeVec
	lynxiBoardPower               gaugeVec
	lynxiDevicePCIEReadBandwidth  gaugeVec
	lynxiDevicePCIEWriteBandwidth gaugeVec
	lynxiDeviceDDRReadBandwidth   gaugeVec
	lynxiDeviceDDRWriteBandwidth  gaugeVec
	smi                           smi.LynSMI
	podRes                        *podresources.PodResources
	devices_cache                 map[int]Props
}

// NewDeviceRecorder 构造一个DeviceRecorder并初始化指标
func NewDeviceRecorder(smi_handle smi.LynSMI, podRes *podresources.PodResources) *DeviceRecorder {
	ret := &DeviceRecorder{
		lynxiDeviceState: newGaugeVec(prometheus.GaugeOpts{
			Name: "lynxi_device_state",
			Help: "The state of the device",
		}, deviceStateLabels),
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
		lynxiDevicePCIEReadBandwidth: newGaugeVec(prometheus.GaugeOpts{
			Name: "lynxi_device_pcie_read_bandwidth",
		}, deviceMetricLabels),
		lynxiDevicePCIEWriteBandwidth: newGaugeVec(prometheus.GaugeOpts{
			Name: "lynxi_device_pcie_write_bandwidth",
		}, deviceMetricLabels),
		lynxiDeviceDDRReadBandwidth: newGaugeVec(prometheus.GaugeOpts{
			Name: "lynxi_device_ddr_read_bandwidth",
		}, deviceMetricLabels),
		lynxiDeviceDDRWriteBandwidth: newGaugeVec(prometheus.GaugeOpts{
			Name: "lynxi_device_ddr_write_bandwidth",
		}, deviceMetricLabels),
		smi:           smi_handle,
		podRes:        podRes,
		devices_cache: make(map[int]Props),
	}
	return ret
}

const (
	labelProductName   = "product_name"
	labelManufacturer  = "manufacturer"
	labelMountTime     = "mount_time"
	labelBoardID       = "board_id"
	labelSerialNumber  = "serial_number"
	labelModel         = "model"
	labelID            = "id"
	labelUUID          = "uuid"
	labelMemTotal      = "mem_total"
	labelPod           = "owner_pod"
	labelContainer     = "owner_container"
	labelNamespace     = "owner_namespace"
	labelErrType       = "err_type"
	labelErrMsg        = "err_msg"
	labelEnableRecover = "enable_recover"
)

type labels []string

var boardLabels labels = labels{labelProductName, labelManufacturer, labelMountTime, labelBoardID, labelSerialNumber}
var deviceMetricLabels labels = append(boardLabels, labels{labelModel, labelID, labelUUID, labelPod, labelNamespace, labelContainer}...)
var deviceStateLabels labels = append(deviceMetricLabels, labels{labelErrType, labelErrMsg, labelEnableRecover}...)
var memLabels labels = append(deviceMetricLabels, labels{labelMemTotal}...)
var ignoreLabels labels = append(boardLabels, labels{labelModel}...)

var mountTime = time.Now().Format(time.RFC3339)

func (ls labels) getValues(device Props) []string {
	ret := make([]string, len(ls))
	for i := range ls {
		ret[i] = ls.getValue(device, i)
	}
	return ret
}

func isLabelIn(label string, labels []string) bool {
	for _, l := range labels {
		if label == l {
			return true
		}
	}
	return false
}

func (ls labels) getNoProps(device Props) []string {
	ret := make([]string, len(ls))
	for i, label := range ls {
		if !isLabelIn(label, ignoreLabels) {
			ret[i] = ls.getValue(device, i)
		}
	}
	return ret
}

var ErrTypName map[smi.ErrType]string = map[smi.ErrType]string{
	smi.ERR_TYPE_CHIP:  "chip",
	smi.ERR_TYPE_BOARD: "board",
	smi.ERR_TYPE_NODE:  "node",
	smi.ERR_TYPE_OTHER: "other",
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
	case labelErrType:
		if device.ErrMsg == nil {
			return ""
		}
		return ErrTypName[device.ErrMsg.Typ]
	case labelErrMsg:
		if device.ErrMsg == nil {
			return ""
		}
		return device.ErrMsg.Msg
	case labelEnableRecover:
		if device.ErrMsg == nil {
			return ""
		}
		return strconv.Itoa(int(device.ErrMsg.EnableRecover))
	default:
		panic("unknown label")
	}
}

func (m *DeviceRecorder) reset() {
	m.lynxiDeviceState.reset()
	m.lynxiDeviceMemUsed.reset()
	m.lynxiDeviceApuUsage.reset()
	m.lynxiDeviceArmUsage.reset()
	m.lynxiDeviceVicUsage.reset()
	m.lynxiDeviceIpeUsage.reset()
	m.lynxiDeviceCurrentTemp.reset()
	m.lynxiBoardPower.reset()
	m.lynxiDevicePCIEReadBandwidth.reset()
	m.lynxiDevicePCIEWriteBandwidth.reset()
	m.lynxiDeviceDDRReadBandwidth.reset()
	m.lynxiDeviceDDRWriteBandwidth.reset()
}

type Props struct {
	podresources.ResourceOwner
	smi.Props
	ID     string
	ErrMsg *smi.ErrMsg
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

func exceptionMapGetOrNil(m map[uint32]smi.ErrMsg, key uint32) *smi.ErrMsg {
	if v, ok := m[key]; ok {
		return &v
	}
	return nil
}

func (m *DeviceRecorder) Record() error {
	var devices smi.PropsMap
	var exceptions = make(map[uint32]smi.ErrMsg)
	var id2ResOwner map[string]podresources.ResourceOwner
	var smiErr error
	var drvErr error
	var resErr error
	concurrentExec(func() {
		devices, smiErr = m.smi.GetDevices()
	}, func() {
		exceptions, drvErr = m.smi.GetDrvExceptionMap()
	}, func() {
		id2ResOwner, resErr = m.getResoureOwners()
	})
	GlobalRecorder.LogIfError(smiErr)
	GlobalRecorder.LogIfError(drvErr)
	if resErr != nil {
		GlobalRecorder.LogError(resErr)
		return resErr
	}
	m.reset()
	for i, props_result := range devices {
		errMsg := exceptionMapGetOrNil(exceptions, uint32(i))
		id := strconv.Itoa(int(i))
		res_owner := id2ResOwner[id]
		props, err := props_result.Get()
		GlobalRecorder.LogIfError(err)
		if props != nil {
			device := Props{res_owner, *props, id, errMsg}
			m.devices_cache[int(i)] = device
			m.lynxiDeviceState.set(device, StateOK)
			m.lynxiDeviceMemUsed.set(device, float64(device.Device.MemoryUsed))
			m.lynxiDeviceApuUsage.set(device, float64(device.Device.ApuUsage))
			m.lynxiDeviceArmUsage.set(device, float64(device.Device.ArmUsage))
			m.lynxiDeviceVicUsage.set(device, float64(device.Device.VicUsage))
			m.lynxiDeviceIpeUsage.set(device, float64(device.Device.IpeUsage))
			m.lynxiDeviceCurrentTemp.set(device, float64(device.Device.Temperature))
			m.lynxiBoardPower.set(device, float64(device.Board.PowerDraw))
			if device.Device.PcieReadBandwidth != nil {
				m.lynxiDevicePCIEReadBandwidth.set(device, float64(*device.Device.PcieReadBandwidth))
			}
			if device.Device.PcieWriteBandwidth != nil {
				m.lynxiDevicePCIEWriteBandwidth.set(device, float64(*device.Device.PcieWriteBandwidth))
			}
			if device.Device.DdrReadBandwidth != nil {
				m.lynxiDeviceDDRReadBandwidth.set(device, float64(*device.Device.DdrReadBandwidth))
			}
			if device.Device.DdrWriteBandwidth != nil {
				m.lynxiDeviceDDRWriteBandwidth.set(device, float64(*device.Device.DdrWriteBandwidth))
			}
		} else {
			props := Props{res_owner, smi.Props{}, id, errMsg}
			props.Device.UUID = m.devices_cache[int(i)].Device.UUID
			m.lynxiDeviceState.setNoProps(props, StateErr)
		}
	}
	return nil
}
