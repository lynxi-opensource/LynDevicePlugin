package lynsmi

import (
	"bufio"
	"encoding/json"
	"net"

	types "lyndeviceplugin/lynsmi-interface"
	"lyndeviceplugin/utils/singleflight"
)

type retType struct {
	props []*types.Props
	err   error
}

var _ types.SMI = &SMIImpl{}

// SMIImpl implements the SMI interface by smiInterface.
type SMIImpl struct {
	reader *bufio.Reader
	conn   net.Conn
	sf     singleflight.Singleflight[retType]
}

func New(addr string) (smi *SMIImpl, err error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return
	}
	return &SMIImpl{bufio.NewReader(conn), conn, singleflight.New[retType]()}, nil
}

func (m *SMIImpl) GetDevices() (types.AllProps, error) {
	ret := m.sf.Fly(func() retType {
		b, err := m.reader.ReadBytes(0)
		if err != nil {
			return retType{nil, err}
		}
		var ret types.AllProps
		err = json.Unmarshal(b[:len(b)-1], &ret)
		return retType{ret, err}
	})
	return ret.props, ret.err
}

func (m *SMIImpl) Close() error {
	return m.conn.Close()
}
