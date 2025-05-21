package main

import (
	"runtime"
	"sync/atomic"
)

type Locker interface {
	Lock()
	UnLock()
}

type SpinLock struct {
	val int32
}

func (s *SpinLock) Lock() {
	for !atomic.CompareAndSwapInt32(&s.val, 0, 1) {
		runtime.Gosched()
	}
}

func (s *SpinLock) UnLock() {
	atomic.StoreInt32(&s.val, 0)
}

func NewSpinLock() Locker {
	return &SpinLock{
		val: 0,
	}
}

func main() {

}
