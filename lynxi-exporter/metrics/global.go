package metrics

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	StateErr float64 = iota
	StateOK
)

var StateEnumDescription string

func (r *globalRecorder) LogIfError(err error) {
	if err != nil {
		r.LogError(err)
	}
}

func (r *globalRecorder) LogError(err error) {
	r.lynxiExporterState.Set(float64(StateErr))
	log.Println(err)
	r.reset <- struct{}{}
}

var GlobalRecorder *globalRecorder

type globalRecorder struct {
	lock               sync.Mutex
	lynxiExporterState prometheus.Gauge
	recoveryDuration   time.Duration
	reset              chan struct{}
}

func init() {
	StateEnumDescription = fmt.Sprintf("%d is ok, %d is err.", int(StateOK), int(StateErr))
	lynxiExporterState := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "lynxi_exporter_state",
		Help: "Is there any error in lynxi-exporter internal, please see the logs, will auto recovering after a while. " + StateEnumDescription,
	})
	lynxiExporterState.Set(float64(StateOK))
	GlobalRecorder = &globalRecorder{
		lynxiExporterState: lynxiExporterState,
		recoveryDuration:   time.Second * 60,
		reset:              make(chan struct{}),
	}
}

func (r *globalRecorder) SetRecoveryDuration(recoveryDuration time.Duration) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.recoveryDuration = recoveryDuration
}

func (r *globalRecorder) Record() error {
	for {
		r.lock.Lock()
		recoveryDuration := r.recoveryDuration
		r.lock.Unlock()
		select {
		case <-time.After(recoveryDuration):
			r.lynxiExporterState.Set(float64(StateOK))
		case <-r.reset:
		}
	}
}
