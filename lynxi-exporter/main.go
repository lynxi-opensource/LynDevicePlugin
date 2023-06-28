package main

import (
	"log"
	"net/http"
	"time"

	smi "lyndeviceplugin/lynsmi-service-client-go"
	"lyndeviceplugin/lynxi-exporter/metrics"

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
	smiC, err := smi.New("127.0.0.1:5432")
	if err != nil {
		log.Fatalln(err)
	}

	// new device recorder and start record
	deviceMetrics := metrics.NewDeviceRecorder(smiC, timeout)
	go func() {
		log.Fatalln(deviceMetrics.Record())
	}()

	podResources := metrics.NewPodContainerRecorder(timeout, smiC)
	go func() {
		log.Fatalln(podResources.Record())
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Println("Listening on :2112. Go to http://localhost:2112/metrics to see metrics.")
	log.Fatalln(http.ListenAndServe(":2112", nil))
}
