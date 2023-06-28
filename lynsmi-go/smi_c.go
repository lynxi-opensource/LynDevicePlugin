package lynsmi

import (
	"fmt"
	"sync"

	"lyndeviceplugin/lynsmi-go/smi_c"
	types "lyndeviceplugin/lynsmi-interface"
	"lyndeviceplugin/utils/errorsext"
	"lyndeviceplugin/utils/singleflight"
)

// var MountTime = time.Now().Format(time.RFC3339)

// Manufacturer is the manufacturer of the board.
const Manufacturer = "lynxi.com"

type retType struct {
	props []*types.Props
	err   error
}

var _ types.SMI = &SMIC{}

// SMIC implements the SMI interface by smiInterface.
type SMIC struct {
	sf singleflight.Singleflight[retType]
}

func New() *SMIC {
	return &SMIC{
		sf: singleflight.New[retType](),
	}
}

func getDevices() (devices []*types.Props, err error) {
	deviceCnt, err := smi_c.GetDeviceCount()
	if err != nil {
		return nil, fmt.Errorf("lynGetDeviceCount() return err: %w", err)
	}
	devices = make([]*types.Props, deviceCnt)
	wg := &sync.WaitGroup{}
	me := errorsext.MultiErrorBuilder{}
	for i := 0; i < int(deviceCnt); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			props, err := smi_c.GetDeviceProperties(int32(i))
			if err == nil {
				devices[i] = &props
			} else {
				me.Push(fmt.Errorf("lynGetdeviceProp(%d) return err: %w", i, err))
			}
		}(i)
	}
	wg.Wait()
	err = me.Error()
	return
}

func (m *SMIC) GetDevices() (types.AllProps, error) {
	ret := m.sf.Fly(func() retType {
		devices, err := getDevices()
		return retType{devices, err}
	})
	return ret.props, ret.err
}

func (m *SMIC) Close() error {
	return nil
}
