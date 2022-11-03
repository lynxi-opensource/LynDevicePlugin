package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"lyndeviceplugin/lynxi-exporter/metrics"
	"lyndeviceplugin/smi"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// TODO: test
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	gr := metrics.GlobalRecorder
	gr.SetRecoveryDuration(time.Second * 60)
	go func() {
		log.Fatalln(gr.Record())
	}()
	timeout := 5 * time.Second
	smiC := smi.NewSMIC()
	deviceID2UUID := make(map[string]string)
	deviceInfos, err := smiC.GetDevices()
	if err != nil {
		log.Fatalln(err)
	}
	for _, deviceInfo := range deviceInfos {
		deviceID2UUID[strconv.Itoa(deviceInfo.ID)] = deviceInfo.UUID
	}

	// new device recorder and start record
	deviceMetrics := metrics.NewDeviceRecorder(smiC, timeout)
	go func() {
		log.Fatalln(deviceMetrics.Record())
	}()

	podResources := metrics.NewPodContainerRecorder(timeout, deviceID2UUID)
	go func() {
		log.Fatalln(podResources.Record())
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Println("Listening on :2112. Go to http://localhost:2112/metrics to see metrics.")
	log.Fatalln(http.ListenAndServe(":2112", nil))
}
