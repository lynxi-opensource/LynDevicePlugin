package main

import (
	"log"
	"time"

	smi "lyndeviceplugin/lynsmi-service-client-go"
	"lyndeviceplugin/lynxi-device-plugin/allocator"
	"lyndeviceplugin/lynxi-device-plugin/server"
	"lyndeviceplugin/lynxi-device-plugin/service"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	smiImpl, err := smi.New("127.0.0.1:5432")
	if err != nil {
		log.Fatalln(err)
	}
	svc := service.NewService(smiImpl, allocator.Alloc{}, time.Second*3)
	s := server.ServerImp{}
	log.Fatalln(s.Run("lynxi_device.sock", "lynxi.com/device", svc))
}
