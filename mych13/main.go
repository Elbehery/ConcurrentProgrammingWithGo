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

type FutexLock struct {
	val int32
}

func (f *FutexLock) Lock() {
	for !atomic.CompareAndSwapInt32(&f.val, 0, 1) {
		futex_wait_int32(&f.val, 1)
	}
}
func (f *FutexLock) UnLock() {
	atomic.StoreInt32(&f.val, 0)
	futex_wake_int32(&f.val, 1)
}

func main() {

}
