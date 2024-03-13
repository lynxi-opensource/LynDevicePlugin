package lynsmi

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLynSMI_GetDevices(t *testing.T) {
	smi := LynSMI{}
	props, err := smi.GetDevices()
	assert.Nil(t, err)
	for id, v := range props {
		fmt.Println(id, v)
	}
}

func TestLynSMI_GetDeviceTopologyList(t *testing.T) {
	smi := LynSMI{}
	list, err := smi.GetDeviceTopologyList()
	assert.Nil(t, err)
	if list == nil {
		fmt.Println("unsupported")
		return
	}
	for _, v := range *list {
		fmt.Println(v)
	}
}
