package smi

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSMIC_GetDevices(t *testing.T) {
	m := NewSMIC()
	got, err := m.GetDevices()
	assert.Nil(t, err)
	assert.NotEmpty(t, got)
	for _, d := range got {
		t.Logf("%+v\n", d)
	}
}

func BenchmarkSMIC_GetDevices(b *testing.B) {
	m := NewSMIC()
	// first, err := m.GetDevices()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := m.GetDevices()
		// assert.Equal(b, first, r)
		assert.Equal(b, err, nil)
	}
}

func BenchmarkSMIC_GetDevicesMultiGoroutines(b *testing.B) {
	m := NewSMIC()
	// first, err := m.GetDevices()
	b.ResetTimer()
	wg := &sync.WaitGroup{}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := m.GetDevices()
			// assert.Equal(b, first, r)
			assert.Equal(b, err, nil)
		}()
	}
	wg.Wait()
}
