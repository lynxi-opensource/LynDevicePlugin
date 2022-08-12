// Package singleflight 用于减少并发获取同一个资源的重复访问
package singleflight

import "sync"

type retTyp[V any] struct {
	v *V
}

// Singleflight 用于并发获取数据时同时只存在一个获取数据的调用
type Singleflight[V any] struct {
	cond *sync.Cond
	ret  *retTyp[V]
}

// New 创建一个 Singleflight
func New[V any]() Singleflight[V] {
	return Singleflight[V]{
		cond: sync.NewCond(&sync.Mutex{}),
		ret:  nil,
	}
}

// Fly 如果已经有获取数据的调用，则返回已经获取的数据，否则通过f获取V，并返回V
func (s *Singleflight[V]) Fly(f func() V) V {
	s.cond.L.Lock()
	var ret = s.ret
	if ret == nil {
		ret = &retTyp[V]{nil}
		s.ret = ret
		s.cond.L.Unlock()
		v := f()
		s.cond.L.Lock()
		ret.v = &v
		s.ret = nil
		s.cond.Broadcast()
		s.cond.L.Unlock()
		return *ret.v
	}
	defer s.cond.L.Unlock()
	for ret.v == nil {
		s.cond.Wait()
	}
	return *ret.v
}
