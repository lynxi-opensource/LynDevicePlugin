package smi

import (
	"fmt"
	"sync"

	"lyndeviceplugin/smi/smi_c"
	"lyndeviceplugin/utils/errorsext"
	"lyndeviceplugin/utils/singleflight"
)

// var MountTime = time.Now().Format(time.RFC3339)

// Manufacturer is the manufacturer of the board.
const Manufacturer = "lynxi.com"

type retType struct {
	boards []Device
	err    error
}

// SMIC implements the SMI interface by smiInterface.
type SMIC struct {
	sf singleflight.Singleflight[retType]
}

func NewSMIC() *SMIC {
	return &SMIC{
		sf: singleflight.New[retType](),
	}
}

func getDevices() (devices []Device, err error) {
	deviceCnt, err := smi_c.GetDeviceCount()
	if err != nil {
		return nil, fmt.Errorf("lynGetDeviceCount() return err: %w", err)
	}
	mtx := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	me := errorsext.MultiErrorBuilder{}
	for i := 0; i < int(deviceCnt); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			deviceProp, err := smi_c.GetDeviceProperties(int32(i))
			mtx.Lock()
			defer mtx.Unlock()
			if err == nil {
				devices = append(devices, Device{
					DeviceInfo: DeviceInfo{
						ID:    i,
						UUID:  deviceProp.UUID,
						Model: deviceProp.Name,
						IsOn:  true,
					},
					BoardInfo: BoardInfo{
						BoardID:      int(deviceProp.BoardID),
						ProductName:  deviceProp.ProductName,
						Manufacturer: Manufacturer,
						// MountTime:    MountTime,
						ChipCnt:      deviceProp.ChipCount,
						SerialNumber: deviceProp.SerialNumber,
					},
					BoardMetrics: BoardMetrics{
						PowerDraw: deviceProp.PowerDraw,
					},
					DeviceMetrics: DeviceMetrics{
						MemUsed:      deviceProp.MemUsed,
						MemTotal:     deviceProp.MemTotal,
						CurrentTemp:  deviceProp.CurrentTemp,
						ApuUsageRate: deviceProp.ApuUsageRate,
						ArmUsageRate: deviceProp.ArmUsageRate,
						VicUsageRate: deviceProp.VicUsageRate,
						IpeUsageRate: deviceProp.IpeUsageRate,
					},
				})
			} else {
				devices = append(devices, Device{DeviceInfo: DeviceInfo{ID: i, IsOn: false}})
				me.Push(fmt.Errorf("lynGetdeviceProp(%d) return err: %w", i, err))
			}
		}(i)
	}
	wg.Wait()
	err = me.Error()
	return
}

func (m *SMIC) GetDevices() ([]Device, error) {
	ret := m.sf.Fly(func() retType {
		devices, err := getDevices()
		return retType{devices, err}
	})
	return ret.boards, ret.err
}
