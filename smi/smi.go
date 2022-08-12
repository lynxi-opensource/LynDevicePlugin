// Package smi 提供获取设备状态和信息的方法
package smi

// Device 描述一个设备的信息和状态
type Device struct {
	BoardInfo
	BoardMetrics
	DeviceInfo
	DeviceMetrics
}

// DeviceInfo 描述设备的基本信息
type DeviceInfo struct {
	ID    int    // 设备的ID
	IsOn  bool   // 设备是否在线
	Model string // 芯片型号，例如KA200
	UUID  string // 设备的UUID
}

// BoardInfo 描述一个板子的信息
type BoardInfo struct {
	ProductName  string // 板卡产品名，例如HP300
	Manufacturer string // 板卡厂家
	// MountTime    string // 板卡接入节点操作系统时间
	BoardID      int    // 板卡ID，例如1
	ChipCnt      uint32 // 芯片数量
	SerialNumber string // 板卡序列号
}

// BoardMetrics 描述一个板子的指标信息
type BoardMetrics struct {
	PowerDraw float32 // 板卡电源消耗, mW
}

// DeviceMetrics 描述一个设备的指标信息
type DeviceMetrics struct {
	MemUsed      uint64 // 内存使用量，单位为字节
	MemTotal     uint64 // 内存总量，单位为字节
	CurrentTemp  int32  // 当前芯片温度，单位为摄氏度
	ApuUsageRate uint32 // APU使用率，单位为百分比
	ArmUsageRate uint32 // ARM使用率，单位为百分比
	VicUsageRate uint32 // VIC使用率，单位为百分比
	IpeUsageRate uint32 // IPE使用率，单位为百分比
}

// SMI 是设备信息提供方必须实现的方法
type SMI interface {
	// GetDevices 获取设备列表
	GetDevices() ([]Device, error)
}
