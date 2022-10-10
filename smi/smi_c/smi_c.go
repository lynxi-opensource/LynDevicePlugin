// Package smi_c 提供smiInterfaceExp.h中定义的函数的go版本
package smi_c

/*
#cgo LDFLAGS: -lLYNSMICLIENTCOMM
#include <lyn_smi.h>
*/
import "C"
import (
	"os"
	"strconv"
)

type Error C.lynError_t

func (e Error) Error() string {
	return strconv.Itoa(int(e))
}

func (e Error) Code() int {
	return int(e)
}

func check(code C.lynError_t) error {
	if code != 0 {
		return Error(code)
	}
	return nil
}

func GetDeviceCount() (int, error) {
	devFiles, err := os.ReadDir("/dev/lynd")
	return len(devFiles) - 1, err
}

type BoardProperties struct {
	ProductName  string
	SerialNumber string
	BoardID      uint32
	ChipCount    uint32
	PowerDraw    float32
}

type DeviceProperties struct {
	BoardProperties
	Name         string
	UUID         string
	MemUsed      uint64
	MemTotal     uint64
	CurrentTemp  int32
	ApuUsageRate uint32
	ArmUsageRate uint32
	VicUsageRate uint32
	IpeUsageRate uint32
}

func newChipProp(raw C.lynDeviceProperties_t) DeviceProperties {
	return DeviceProperties{
		BoardProperties: BoardProperties{
			ProductName:  C.GoString(&raw.boardProductName[0]),
			SerialNumber: C.GoString(&raw.boardSerialNumber[0]),
			BoardID:      uint32(raw.boardId),
			ChipCount:    uint32(raw.boardChipCount),
			PowerDraw:    float32(raw.boardPowerDraw),
		},

		Name:         C.GoString(&raw.deviceName[0]),
		UUID:         C.GoString(&raw.deviceUuid[0]),
		MemUsed:      uint64(raw.deviceMemoryUsed),
		MemTotal:     uint64(raw.deviceMemoryTotal),
		CurrentTemp:  int32(raw.deviceTemperatureCurrent),
		ApuUsageRate: uint32(raw.deviceApuUsageRate),
		ArmUsageRate: uint32(raw.deviceArmUsageRate),
		VicUsageRate: uint32(raw.deviceVicUsageRate),
		IpeUsageRate: uint32(raw.deviceIpeUsageRate),
	}
}

func GetDeviceProperties(devID int32) (DeviceProperties, error) {
	var raw C.lynDeviceProperties_t
	err := check(C.lynGetDeviceProperties(C.int32_t(devID), &raw))
	return newChipProp(raw), err
}
