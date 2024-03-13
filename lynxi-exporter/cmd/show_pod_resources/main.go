package main

import (
	"fmt"
	podresources "lyndeviceplugin/lynxi-exporter/pod_resources"
)

func main() {
	m, err := podresources.New()
	if err != nil {
		fmt.Println(err)
		return
	}
	ret, err := m.Get()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ret)
}
