package main

import (
	"log"
	"time"

	"lyndeviceplugin/lynsmi-service-client-go"
	"lyndeviceplugin/lynxi-device-plugin/allocator"
	"lyndeviceplugin/lynxi-device-plugin/server"
	"lyndeviceplugin/lynxi-device-plugin/service"
)

type AllocationType int

const (
	ALLOCATION_TYPE_DEVICE AllocationType = iota
	ALLOCATION_TYPE_HP280
)

func getAllocationType() AllocationType {
	smi := lynsmi.LynSMI{}
	devices, err := smi.GetDevices()
	if err != nil {
		log.Fatalln("smi GetDevices err: ", err)
	}
	hasHP280 := false
	onlyHP280 := true
	for _, d := range devices {
		if v, err := d.Get(); err == nil && v != nil {
			if v.Board.ProductName == service.HP280Name {
				hasHP280 = true
			} else {
				onlyHP280 = false
			}
		}
	}
	if hasHP280 && onlyHP280 {
		return ALLOCATION_TYPE_HP280
	}
	return ALLOCATION_TYPE_DEVICE
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	allocationType := getAllocationType()
	log.Println("allocationType", allocationType)

	go func() {
		var svc service.Service
		switch allocationType {
		case ALLOCATION_TYPE_DEVICE:
			svc = service.NewDeviceAllocationService(allocator.Alloc{}, time.Second*3)
		case ALLOCATION_TYPE_HP280:
			svc = service.NewHP280AllocationService(allocator.Alloc{}, time.Second*3)
		default:
			panic("unreachable")
		}
		s := server.ServerImp{}
		log.Fatalln(s.Run("lynxi_device.sock", svc))
	}()

	ticker := time.NewTicker(time.Second * 60)
	for range ticker.C {
		if getAllocationType() != allocationType {
			log.Println("allocation type changed")
			return
		}
	}
}
