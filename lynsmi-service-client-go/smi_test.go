package lynsmi

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSMIImpl_GetDevices(t *testing.T) {
	smi, err := New("127.0.0.1:5432")
	assert.Nil(t, err)
	props, err := smi.GetDevices()
	assert.Nil(t, err)
	for _, v := range props {
		fmt.Println(v)
	}
}
