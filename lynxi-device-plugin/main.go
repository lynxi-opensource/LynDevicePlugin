package main

import (
	"log"
	"time"

	"lyndeviceplugin/lynxi-device-plugin/allocator"
	"lyndeviceplugin/lynxi-device-plugin/server"
	"lyndeviceplugin/lynxi-device-plugin/service"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	svc := service.NewService(allocator.Alloc{}, time.Second*3)
	s := server.ServerImp{}
	log.Fatalln(s.Run("lynxi_device.sock", "lynxi.com/device", svc))
}
