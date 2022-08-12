package smi_c

import (
	"bytes"
	"os/exec"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDeviceCount(t *testing.T) {
	got, err := GetDeviceCount()
	assert.Nil(t, err)
	tmp, err := exec.Command("sh", "-c", "ls /dev/lynd/ | wc -w").Output()
	assert.Nil(t, err)
	want, err := strconv.Atoi(string(bytes.TrimSpace(tmp)))
	want -= 1 // exclude the manager device
	assert.Nil(t, err)
	assert.Equal(t, want, got)
}

func TestGetDeviceProperties(t *testing.T) {
	cnt, err := GetDeviceCount()
	assert.Nil(t, err)
	for i := 0; i < cnt; i++ {
		got, err := GetDeviceProperties(int32(i))
		assert.Nil(t, err)
		t.Logf("%+v\n", got)
	}
}
