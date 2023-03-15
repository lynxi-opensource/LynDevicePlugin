package metrics

import (
	"strconv"
	"time"

	"lyndeviceplugin/smi"

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

func (m gaugeVec) set(device smi.Device, v float64) {
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
	timeout                time.Duration
}

// NewDeviceRecorder 构造一个DeviceRecorder并初始化指标
func NewDeviceRecorder(smi smi.SMI, timeout time.Duration) *DeviceRecorder {
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
		smi:     smi,
		timeout: timeout,
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
)

type labels []string

var boardLabels labels = labels{labelProductName, labelManufacturer, labelMountTime, labelBoardID, labelSerialNumber}
var deviceMetricLabels labels = append(boardLabels, labels{labelModel, labelID, labelUUID}...)
var memLabels labels = append(deviceMetricLabels, labels{labelMemTotal}...)

var mountTime = time.Now().Format(time.RFC3339)

func (ls labels) getValues(device smi.Device) []string {
	ret := make([]string, len(ls))
	for i := range ls {
		ret[i] = ls.getValue(device, i)
	}
	return ret
}

func (ls labels) getValue(device smi.Device, i int) string {
	switch ls[i] {
	case labelProductName:
		return device.ProductName
	case labelManufacturer:
		return device.Manufacturer
	case labelMountTime:
		return mountTime
	case labelBoardID:
		return strconv.Itoa(device.BoardID)
	case labelSerialNumber:
		return device.SerialNumber
	case labelModel:
		return device.Model
	case labelID:
		return strconv.Itoa(device.ID)
	case labelUUID:
		return device.UUID
	case labelMemTotal:
		return strconv.Itoa(int(device.MemTotal))
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

// Record 一直阻塞不会返回错误，外部通过lynxi_exporter_state或日志查看exporter的状态是否正常
func (m *DeviceRecorder) Record() error {
	ticker := time.NewTicker(m.timeout)
	for range ticker.C {
		devices, err := m.smi.GetDevices()
		GlobalRecorder.logIfError(err)
		m.reset()
		for _, device := range devices {
			if device.IsOn {
				m.lynxiDeviceStates.set(device, StateOK)
				m.lynxiDeviceMemUsed.set(device, float64(device.MemUsed))
				m.lynxiDeviceApuUsage.set(device, float64(device.ApuUsageRate))
				m.lynxiDeviceArmUsage.set(device, float64(device.ArmUsageRate))
				m.lynxiDeviceVicUsage.set(device, float64(device.VicUsageRate))
				m.lynxiDeviceIpeUsage.set(device, float64(device.IpeUsageRate))
				m.lynxiDeviceCurrentTemp.set(device, float64(device.CurrentTemp))
				m.lynxiBoardPower.set(device, float64(device.PowerDraw))
			} else {
				m.lynxiDeviceStates.set(device, StateErr)
			}
		}
	}
	return nil
}
