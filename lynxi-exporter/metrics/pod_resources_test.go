package metrics

import (
	"math/rand"
	"sort"
	"strconv"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/stretchr/testify/assert"
)

func TestNumberStringSlice(t *testing.T) {
	for range [100]struct{}{} {
		a := make([]int, 100)
		s := make([]string, 100)
		for i := range a {
			a[i] = rand.Intn(100)
			s[i] = strconv.Itoa(a[i])
		}
		sort.Sort(sort.IntSlice(a))
		sort.Sort(StringNumberSlice(s))
		for i := range a {
			assert.Equal(t, s[i], strconv.Itoa(a[i]))
		}
	}
}

func TestLabelLengthLimit(t *testing.T) {
	g := gaugeVec{
		m: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "test",
			Help: "test",
		}, []string{"test"})}
	g.m.WithLabelValues(string(make([]byte, 1000000))).Set(1)
	// http.Handle("/metrics", promhttp.Handler())
	// log.Println("Listening on :2112. Go to http://localhost:2112/metrics to see metrics.")
	// log.Fatalln(http.ListenAndServe(":2112", nil))
}
