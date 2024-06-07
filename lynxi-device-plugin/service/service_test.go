package service

import (
	"fmt"
	"lyndeviceplugin/lynsmi-service-client-go"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestForEachDeviceList(t *testing.T) {
	scoreMap, err := getScoreMap(lynsmi.LynSMI{})
	assert.Nil(t, err)
	// fmt.Println(scoreMap)

	must := []int{}
	availible := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	availible = sliceSub(availible, must)
	forEachDeviceList(scoreMap, 0, 0, must, availible, 4, func(score int32, dist int32, selected []int) bool {
		fmt.Println(score, dist, selected)
		return true
	})
}

func TestGetBestDeviceList(t *testing.T) {
	scoreMap, err := getScoreMap(lynsmi.LynSMI{})
	assert.Nil(t, err)

	must := []int{}
	availible := []int{0, 1, 2, 3, 4, 5, 6, 7, 8}
	availible = sliceSub(availible, must)

	ret := getBestDeviceList(scoreMap, must, availible, 2)
	fmt.Println(ret)
}

func TestGetBestDeviceList_64Device(t *testing.T) {
	scoreMap := make(map[[2]int][2]int32)
	must := []int{}
	availible := make([]int, 64)
	for i := range availible {
		availible[i] = i
		for j := range availible {
			if i != j {
				scoreMap[[2]int{i, j}] = [2]int32{0, 0}
			}
		}
	}
	ret := getBestDeviceList(scoreMap, must, availible, 6)
	fmt.Println(ret)
}

func TestGetDeviceListByGreedyPolicy_64Device(t *testing.T) {
	scoreMap := make(map[[2]int][2]int32)
	must := []int{}
	availible := make([]int, 64)
	for i := range availible {
		availible[i] = i
		for j := range availible {
			if i != j {
				scoreMap[[2]int{i, j}] = [2]int32{int32(rand.Intn(100)), int32(rand.Intn(100))}
			}
		}
	}
	ret := getDeviceListByGreedyPolicy(scoreMap, must, availible, 6)
	fmt.Println(ret)
	ret = getBestDeviceList(scoreMap, must, availible, 6)
	fmt.Println(ret)
}

func TestGetDeviceListByGreedyPolicy(t *testing.T) {
	scoreMap, err := getScoreMap(lynsmi.LynSMI{})
	assert.Nil(t, err)

	must := []int{}
	availible := []int{0, 1, 2, 3, 4, 5, 6, 7, 8}
	availible = sliceSub(availible, must)

	ret := getDeviceListByGreedyPolicy(scoreMap, must, availible, 7)
	fmt.Println(ret)
}
