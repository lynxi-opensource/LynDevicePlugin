package lynsmi

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSMIImpl_GetDevices(t *testing.T) {
	smi, err := New("127.0.0.1:5432")
	assert.Nil(t, err)
	props, err := smi.GetDevices()
	assert.Nil(t, err)
	for _, v := range props {
		fmt.Println(v)
	}
}

func TestSMIImpl_GetDevices_Concurrency(t *testing.T) {
	smi, err := New("127.0.0.1:5432")
	assert.Nil(t, err)
	wg := &sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := smi.GetDevices()
			assert.Nil(t, err)
		}()
	}
	wg.Wait()
}

func BenchmarkSMIImpl_GetDevices(b *testing.B) {
	smi, err := New("127.0.0.1:5432")
	assert.Nil(b, err)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := smi.GetDevices()
		assert.Nil(b, err)
	}
}
