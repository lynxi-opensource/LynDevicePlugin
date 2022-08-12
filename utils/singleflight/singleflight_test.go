package singleflight

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestSingleflight(t *testing.T) {
	runtime.GOMAXPROCS(1)
	// f, _ := os.Create("cpu.out")
	// pprof.StartCPUProfile(f)
	// defer pprof.StopCPUProfile()
	sf := New[int32]()
	for x := range [5]struct{}{} {
		fmt.Println("x = ", x)
		wg := sync.WaitGroup{}
		for i := range [1000000]struct{}{} {
			wg.Add(1)
			go func(i, x int32) {
				defer wg.Done()
				if sf.Fly(func() int32 {
					fmt.Printf("do work with i = %d and x = %d\n", i, x)
					time.Sleep(time.Millisecond * 200)
					return x
				}) != x {
					t.Errorf("not equal")
				}
			}(int32(i), int32(x))
		}
		wg.Wait()
	}
}
