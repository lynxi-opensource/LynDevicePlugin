package lynsmi

import (
	"errors"
	"io"
	"net/http"
	"time"

	"k8s.io/apimachinery/pkg/util/json"
)

type PropsMap map[int32]Result[Props]
type P2PAttrList []Result[DeviceP2PAttr]

// Props 描述一个设备的信息和状态
type Props struct {
	Board  BoardProps  `json:"board"`
	Device DeviceProps `json:"device"`
}

// DeviceProps 描述设备的基本信息
type DeviceProps struct {
	Name        string `json:"name"`         // 芯片型号，例如KA200
	UUID        string `json:"uuid"`         // 设备的UUID
	MemoryUsed  uint64 `json:"memory_used"`  // 内存使用量，单位为字节
	MemoryTotal uint64 `json:"memory_total"` // 内存总量，单位为字节
	Temperature int32  `json:"temperature"`  // 当前芯片温度，单位为摄氏度
	ApuUsage    uint32 `json:"apu_usage"`    // APU使用率，单位为百分比
	ArmUsage    uint32 `json:"arm_usage"`    // ARM使用率，单位为百分比
	VicUsage    uint32 `json:"vic_usage"`    // VIC使用率，单位为百分比
	IpeUsage    uint32 `json:"ipe_usage"`    // IPE使用率，单位为百分比
}

// BoardProps 描述一个板子的信息
type BoardProps struct {
	ProductName  string  `json:"product_name"`  // 板卡产品名，例如HP300
	Brand        string  `json:"brand"`         // 板卡厂家
	SerialNumber string  `json:"serial_number"` // 板卡序列号
	ID           uint32  `json:"id"`            // 板卡ID，例如1
	ChipCount    uint32  `json:"chip_count"`    // 芯片数量
	PowerDraw    float32 `json:"power_draw"`    // 板卡电源消耗, mW
}

type P2PMode string

const (
	NonSupport P2PMode = "NonSupport"
	P2PLinkPIX P2PMode = "P2PLinkPIX"
	P2PLinkPXB P2PMode = "P2PLinkPXB"
	P2PLinkPHB P2PMode = "P2PLinkPHB"
	P2PLinkSYS P2PMode = "P2PLinkSYS"
)

type P2PAttr struct {
	Mode P2PMode `json:"mode"`
	Dist int32   `json:"dist"`
}

type DeviceP2PAttr struct {
	DevicePair [2]int32 `json:"device_pair"`
	Attr       P2PAttr  `json:"attr"`
}

type Result[T any] struct {
	Ok  *T      `json:"Ok,omitempty"`
	Err *string `json:"Err,omitempty"`
}

func (r Result[T]) Get() (*T, error) {
	if r.Err != nil {
		return nil, errors.New(*r.Err)
	}
	return r.Ok, nil
}

// LynSMI implements the SMI interface by smiInterface.
type LynSMI struct{}

func (m LynSMI) GetDevices() (ret PropsMap, err error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get("http://localhost:5432/devices")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &ret)
	return
}

func (m LynSMI) GetDeviceTopologyList() (ret *P2PAttrList, err error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get("http://localhost:5432/device_topology_list")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &ret)
	return
}
