package main

import (
	"log"
	"time"

	"lyndeviceplugin/lynxi-k8s-device-plugin/allocator"
	"lyndeviceplugin/lynxi-k8s-device-plugin/server"
	"lyndeviceplugin/lynxi-k8s-device-plugin/service"
	"lyndeviceplugin/smi"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	smiImpl := smi.NewSMIC()
	crash := make(chan error)
	svc := service.NewService(smiImpl, allocator.Alloc{}, crash, time.Second*3)
	s := server.ServerImp{Crash: crash}
	log.Fatalln(s.Run("lynxi_device.sock", "lynxi.com/device", svc))
}
