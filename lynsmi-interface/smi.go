// Package smi 提供获取设备状态和信息的方法
package lynsmi_interface

type AllProps []*Props

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

// SMI 是设备信息提供方必须实现的方法
type SMI interface {
	// GetDevices 获取设备列表
	Close() error
	GetDevices() (AllProps, error)
}
