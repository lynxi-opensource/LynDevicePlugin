// Package smi_c 提供smiInterfaceExp.h中定义的函数的go版本
package smi_c

import (
	"context"
	"errors"
	"fmt"
	types "lyndeviceplugin/lynsmi-interface"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"time"
	"unsafe"
)

/*
#cgo LDFLAGS: -lLYNSMICLIENTCOMM
#include <lyn_smi.h>

typedef struct
{
    char boardProductName[ARRAY_MAX_LEN];
    char boardBrand[ARRAY_MAX_LEN];
    char boardFirmwareVersion[ARRAY_MAX_LEN];
    char boardProductNumber[ARRAY_MAX_LEN];
    char boardSerialNumber[ARRAY_MAX_LEN];
    uint32_t boardId;
    uint32_t boardChipCount;
    float boardPowerDraw;
    float boardPowerLimit;
    float boardVoltage;

    char deviceName[ARRAY_MAX_LEN];
    char deviceUuid[ARRAY_MAX_LEN];
    uint64_t deviceApuClockFrequency;
    uint64_t deviceApuClockFrequencyLimit;
    uint64_t deviceArmClockFrequency;
    uint64_t deviceArmClockFrequencyLimit;
    uint64_t deviceMemClockFrequency;
    uint64_t deviceMemClockFrequencyLimit;
    uint64_t deviceMemoryUsed;
    uint64_t deviceMemoryTotal;
    int32_t deviceTemperatureCurrent;
    int32_t deviceTemperatureSlowdown;
    int32_t deviceTemperatureLimit;
    uint32_t deviceApuUsageRate;
    uint32_t deviceArmUsageRate;
    uint32_t deviceVicUsageRate;
    uint32_t deviceIpeUsageRate;
    uint32_t deviceEccStat;
    uint32_t deviceDdrErrorCount;
    uint32_t deviceDdrNoErrorCount;
    float deviceVoltage;

    uint32_t processCount;
    uint32_t pid[PROCESS_COUNT_LIMIT];
    char processName[PROCESS_COUNT_LIMIT][PROCESS_NAME_LEN];
    uint64_t processUseMemory[PROCESS_COUNT_LIMIT];
} lynDevicePropertiesOld_t;
*/
import "C"

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

func GetDeviceCountByDir() (int, error) {
	devFiles, err := os.ReadDir("/dev/lynd")
	cnt := 0
	for _, f := range devFiles {
		if _, err := strconv.ParseInt(f.Name(), 10, 64); err == nil {
			cnt++
		}
	}
	return cnt, err
}

type DriverVersion [3]int

func (v DriverVersion) String() string {
	return fmt.Sprintf("%d.%d.%d", v[0], v[1], v[2])
}

func (v DriverVersion) LessThan(other DriverVersion) bool {
	for i := range v {
		if v[i] == other[i] {
			continue
		}
		return v[i] < other[i]
	}
	return false
}

func NewDriverVersionFromBytes(bytes []byte) (ret DriverVersion, err error) {
	r, err := regexp.Compile(`(\d+)\.(\d+)\.(\d+)`)
	if err != nil {
		return
	}
	matches := r.FindSubmatch(bytes)
	if len(matches) != 4 {
		err = errors.New("match version number failed from: " + string(bytes))
		return
	}
	for i, v := range matches[1:] {
		ret[i], err = strconv.Atoi(string(v))
		if err != nil {
			return
		}
	}
	return
}

func NewDriverVersionFromSMIBin() (ret DriverVersion, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	output, err := exec.CommandContext(ctx, "lynxi-smi", "-v").Output()
	if err != nil {
		return
	}
	return NewDriverVersionFromBytes(output)
}

var useOldStructBefore DriverVersion = DriverVersion{1, 10, 2}

var isUseOldStruct bool

func init() {
	current_version, err := NewDriverVersionFromSMIBin()
	if err != nil {
		panic(err)
	}
	isUseOldStruct = current_version.LessThan(useOldStructBefore)
}

func GetDeviceCount() (int, error) {
	var ret C.int32_t
	err := check(C.lynGetDeviceCountSmi(&ret))
	return int(ret), err
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

func newChipProp(raw C.lynDeviceProperties_t) types.Props {
	return types.Props{
		Board: types.BoardProps{
			ProductName:  C.GoString(&raw.boardProductName[0]),
			SerialNumber: C.GoString(&raw.boardSerialNumber[0]),
			Brand:        C.GoString(&raw.boardBrand[0]),
			ID:           uint32(raw.boardId),
			ChipCount:    uint32(raw.boardChipCount),
			PowerDraw:    float32(raw.boardPowerDraw),
		},
		Device: types.DeviceProps{
			Name:        C.GoString(&raw.deviceName[0]),
			UUID:        C.GoString(&raw.deviceUuid[0]),
			MemoryUsed:  uint64(raw.deviceMemoryUsed),
			MemoryTotal: uint64(raw.deviceMemoryTotal),
			Temperature: int32(raw.deviceTemperatureCurrent),
			ApuUsage:    uint32(raw.deviceApuUsageRate),
			ArmUsage:    uint32(raw.deviceArmUsageRate),
			VicUsage:    uint32(raw.deviceVicUsageRate),
			IpeUsage:    uint32(raw.deviceIpeUsageRate),
		},
	}
}

func newChipPropOld(raw C.lynDevicePropertiesOld_t) types.Props {
	return types.Props{
		Board: types.BoardProps{
			ProductName:  C.GoString(&raw.boardProductName[0]),
			SerialNumber: C.GoString(&raw.boardSerialNumber[0]),
			Brand:        C.GoString(&raw.boardBrand[0]),
			ID:           uint32(raw.boardId),
			ChipCount:    uint32(raw.boardChipCount),
			PowerDraw:    float32(raw.boardPowerDraw),
		},
		Device: types.DeviceProps{
			Name:        C.GoString(&raw.deviceName[0]),
			UUID:        C.GoString(&raw.deviceUuid[0]),
			MemoryUsed:  uint64(raw.deviceMemoryUsed),
			MemoryTotal: uint64(raw.deviceMemoryTotal),
			Temperature: int32(raw.deviceTemperatureCurrent),
			ApuUsage:    uint32(raw.deviceApuUsageRate),
			ArmUsage:    uint32(raw.deviceArmUsageRate),
			VicUsage:    uint32(raw.deviceVicUsageRate),
			IpeUsage:    uint32(raw.deviceIpeUsageRate),
		},
	}
}

func GetDevicePropertiesNew(devID int32) (types.Props, error) {
	var raw C.lynDeviceProperties_t
	err := check(C.lynGetDeviceProperties(C.int32_t(devID), &raw))
	return newChipProp(raw), err
}

func GetDevicePropertiesOld(devID int32) (types.Props, error) {
	var raw C.lynDevicePropertiesOld_t
	err := check(C.lynGetDeviceProperties(C.int32_t(devID), (*C.lynDeviceProperties_t)(unsafe.Pointer(&raw))))
	return newChipPropOld(raw), err
}

func GetDeviceProperties(devID int32) (types.Props, error) {
	if isUseOldStruct {
		return GetDevicePropertiesOld(devID)
	}
	return GetDevicePropertiesNew(devID)
}
