package lynsmi

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"

	types "lyndeviceplugin/lynsmi-interface"
)

var _ types.SMI = &SMIImpl{}

// SMIImpl implements the SMI interface by smiInterface.
type SMIImpl struct {
	reader *bufio.Reader
	conn   net.Conn
}

type InitError string

func (m InitError) Error() string {
	return string(m)
}

func New(addr string) (smi *SMIImpl, err error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return
	}
	return &SMIImpl{bufio.NewReader(conn), conn}, nil
}

func (m *SMIImpl) GetDevices() (ret types.AllProps, err error) {
	b, err := m.reader.ReadBytes(0)
	if err != nil {
		return
	}
	fmt.Println(string(b))
	err = json.Unmarshal(b[:len(b)-1], &ret)
	return
}

func (m *SMIImpl) Close() error {
	return m.conn.Close()
}
