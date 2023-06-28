package smi_c

import (
	"bytes"
	"os/exec"
	"reflect"
	"strconv"
	"testing"

	"fmt"

	"github.com/stretchr/testify/assert"
)

func TestGetDeviceCount(t *testing.T) {
	got, err := GetDeviceCount()
	assert.Nil(t, err)
	tmp, err := exec.Command("sh", "-c", "ls /dev/lynd/ | wc -w").Output()
	assert.Nil(t, err)
	want, err := strconv.Atoi(string(bytes.TrimSpace(tmp)))
	want -= 2 // exclude the manager device
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

func TestNewDriverVersionFromBytes(t *testing.T) {
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name    string
		args    args
		wantRet DriverVersion
		wantErr bool
	}{
		{"normal1", args{[]byte("SMI version: 1.9.0")}, DriverVersion{1, 9, 0}, false},
		{"normal2", args{[]byte("SMI version: 1.11.0")}, DriverVersion{1, 11, 0}, false},
		{"normal3", args{[]byte("SMI version: 2.9.0")}, DriverVersion{2, 9, 0}, false},
		{"normal4", args{[]byte("SMI version: 0.9.0")}, DriverVersion{0, 9, 0}, false},
		{"normal5", args{[]byte("SMI version: 1.9.10")}, DriverVersion{1, 9, 10}, false},
		{"error", args{[]byte("SMI version: ")}, DriverVersion{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRet, err := NewDriverVersionFromBytes(tt.args.bytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDriverVersionFromBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRet, tt.wantRet) {
				t.Errorf("NewDriverVersionFromBytes() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}

func TestDriverVersion_LessThan(t *testing.T) {
	type args struct {
		other DriverVersion
	}
	tests := []struct {
		name string
		v    DriverVersion
		args args
		want bool
	}{
		{"less than", DriverVersion{0, 9, 0}, args{DriverVersion{1, 11, 0}}, true},
		{"less than", DriverVersion{1, 9, 0}, args{DriverVersion{1, 11, 0}}, true},
		{"less than", DriverVersion{1, 10, 1}, args{DriverVersion{1, 11, 0}}, true},
		{"equal", DriverVersion{1, 11, 0}, args{DriverVersion{1, 11, 0}}, false},
		{"greater than", DriverVersion{1, 11, 1}, args{DriverVersion{1, 11, 0}}, false},
		{"greater than", DriverVersion{1, 12, 0}, args{DriverVersion{1, 11, 0}}, false},
		{"greater than", DriverVersion{2, 1, 0}, args{DriverVersion{1, 11, 0}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.LessThan(tt.args.other); got != tt.want {
				t.Errorf("DriverVersion.LessThan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDriverVersionFromSMIBin(t *testing.T) {
	tests := []struct {
		name    string
		wantRet DriverVersion
		wantErr bool
	}{
		{"test", DriverVersion{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRet, err := NewDriverVersionFromSMIBin()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDriverVersionFromSMIBin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(gotRet)
		})
	}
}
