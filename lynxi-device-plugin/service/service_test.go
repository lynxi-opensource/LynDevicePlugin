package service

import (
	"fmt"
	"lyndeviceplugin/lynsmi-service-client-go"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestForEachDeviceList(t *testing.T) {
	must := []string{"1", "0"}
	availible := []string{"0", "1", "2", "3", "4", "5"}
	availible = sliceSub(availible, must)
	forEachDeviceList(must, availible, 5, func(selected []string) bool {
		fmt.Println(selected)
		return true
	})
}

func TestGetBestDeviceList(t *testing.T) {
	scoreMap, err := getScoreMap(lynsmi.LynSMI{})
	assert.Nil(t, err)

	must := []string{}
	availible := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8"}
	availible = sliceSub(availible, must)

	ret := getBestDeviceList(scoreMap, must, availible, 2)
	fmt.Println(ret)
}
