package main

import (
	"log"
	"net/http"
	"time"

	smi "lyndeviceplugin/lynsmi-service-client-go"
	"lyndeviceplugin/lynxi-exporter/metrics"
	podresources "lyndeviceplugin/lynxi-exporter/pod_resources"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func startRecord(recorder metrics.Recorder, tickTime time.Duration) {
	go func() {
		ticker := time.NewTicker(tickTime)
		for range ticker.C {
			if err := recorder.Record(); err != nil {
				log.Fatalln(err)
			}
		}
	}()
}

// TODO: test
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	gr := metrics.GlobalRecorder
	gr.SetRecoveryDuration(time.Second * 60)
	go func() {
		log.Fatalln(gr.Record())
	}()
	smiImpl, err := smi.New("127.0.0.1:5432")
	if err != nil {
		log.Fatalln(err)
	}
	defer smiImpl.Close()

	podRes, err := podresources.New()
	if err != nil {
		log.Fatalln(err)
	}
	defer podRes.Close()
	log.Println("connect kubelet")

	tickTime := 1 * time.Second

	// new device recorder and start record
	deviceMetrics := metrics.NewDeviceRecorder(smiImpl, podRes)
	startRecord(deviceMetrics, tickTime)

	podResources := metrics.NewPodContainerRecorder(smiImpl, podRes)
	startRecord(podResources, tickTime)

	http.Handle("/metrics", promhttp.Handler())
	log.Println("Listening on :2112. Go to http://localhost:2112/metrics to see metrics.")
	log.Fatalln(http.ListenAndServe(":2112", nil))
}
