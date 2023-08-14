package main

import (
	"fmt"
	"sync"
)

type WaitGrp struct {
	size int
	cond *sync.Cond
}

func NewWaitGrp() *WaitGrp {
	return &WaitGrp{cond: sync.NewCond(&sync.Mutex{})}
}

func (wg *WaitGrp) Add(delta int) {
	wg.cond.L.Lock()
	wg.size += delta
	wg.cond.L.Unlock()
}

func (wg *WaitGrp) Wait() {
	wg.cond.L.Lock()
	for wg.size > 0 {
		wg.cond.Wait()
	}
	wg.cond.L.Unlock()
}

func (wg *WaitGrp) Done() {
	wg.cond.L.Lock()
	wg.size--
	if wg.size == 0 {
		wg.cond.Broadcast()
	}
	wg.cond.L.Unlock()
}

// MyWaitGroup is a wait-group implementation based on Semaphore.
type MyWaitGroup struct {
	s *sema
}

func NewMyWaitGroup(n int) *MyWaitGroup {
	return &MyWaitGroup{newSema(1 - n)}
}

func (wg *MyWaitGroup) Wait() {
	wg.s.acquire()
}

func (wg *MyWaitGroup) Done() {
	wg.s.release()
}

// sema is to be used internally by MyWaitGroup type.
type sema struct {
	permits int
	cond    *sync.Cond
}

func newSema(n int) *sema {
	return &sema{
		permits: n,
		cond:    sync.NewCond(&sync.Mutex{}),
	}
}

func (s *sema) acquire() {
	s.cond.L.Lock()
	for s.permits <= 0 {
		s.cond.Wait()
	}
	s.permits--
	s.cond.L.Unlock()
}

func (s *sema) release() {
	s.cond.L.Lock()
	s.permits++
	s.cond.Signal()
	s.cond.L.Unlock()
}

func main() {
	wg := NewMyWaitGroup(4)
	for i := 1; i <= 4; i++ {
		go doWork(i, wg)
	}
	wg.Wait()
	fmt.Println("All complete")
}

func doWork(id int, wg *MyWaitGroup) {
	fmt.Println(id, "Done working ")
	wg.Done()
}
