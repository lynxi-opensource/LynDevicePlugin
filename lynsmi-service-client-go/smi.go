package lynsmi

import (
	"bufio"
	"encoding/json"
	"net"
	"sync"

	types "lyndeviceplugin/lynsmi-interface"
)

var _ types.SMI = &SMIImpl{}

type propsWithID struct {
	ID    int          `json:"id"`
	Props *types.Props `json:"props"`
	Err   *string      `json:"err"`
}

// SMIImpl implements the SMI interface by smiInterface.
type SMIImpl struct {
	conn     net.Conn
	allProps types.AllProps
	mtx      *sync.Mutex
	cond     *sync.Cond
	err      error
}

func New(addr string, devicesErrorHandler func(int, string)) (smi *SMIImpl, err error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return
	}
	mtx := &sync.Mutex{}
	cond := sync.NewCond(mtx)
	smi = &SMIImpl{conn, make(types.AllProps, 0), mtx, cond, nil}
	go func() {
		reader := bufio.NewReader(conn)
		for {
			b, err := reader.ReadBytes(0)
			if err != nil {
				smi.err = err
				return
			}
			var ret propsWithID
			err = json.Unmarshal(b[:len(b)-1], &ret)
			if err != nil {
				smi.err = err
				return
			}
			mtx.Lock()
			if ret.ID >= len(smi.allProps) {
				smi.allProps = append(smi.allProps, make(types.AllProps, ret.ID+1-len(smi.allProps))...)
			}
			smi.allProps[ret.ID] = ret.Props
			if ret.Err != nil {
				devicesErrorHandler(ret.ID, *ret.Err)
			}
			smi.cond.Broadcast()
			mtx.Unlock()
		}
	}()
	return smi, nil
}

func (m *SMIImpl) GetDevices() (types.AllProps, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.cond.Wait()
	return m.allProps, m.err
}

func (m *SMIImpl) Close() error {
	return m.conn.Close()
}
