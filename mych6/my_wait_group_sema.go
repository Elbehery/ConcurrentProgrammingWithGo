package main

import (
	"fmt"
	"sync"
)

type semaphore struct {
	permits int
	cnd     *sync.Cond
}

func newSemaphore(n int) *semaphore {
	return &semaphore{
		permits: n,
		cnd:     sync.NewCond(&sync.Mutex{}),
	}
}

func (s *semaphore) acquire() {
	s.cnd.L.Lock()
	for s.permits <= 0 {
		s.cnd.Wait()
	}
	s.cnd.L.Unlock()
}

func (s *semaphore) release() {
	s.cnd.L.Lock()
	s.permits++
	s.cnd.Signal()
	s.cnd.L.Unlock()
}

type MyWaitGroupSema struct {
	sema *semaphore
}

func NewMyWaitGroupSema(size int) *MyWaitGroupSema {
	return &MyWaitGroupSema{sema: newSemaphore(1 - size)}
}

func (wg *MyWaitGroupSema) Wait() {
	wg.sema.acquire()
}

func (wg *MyWaitGroupSema) Done() {
	wg.sema.release()
}

func doWork(id int, wg *MyWaitGroupSema) {
	fmt.Println(id, "Done working ")
	wg.Done()
}

func main() {
	wg := NewMyWaitGroupSema(4)
	for i := 1; i <= 4; i++ {
		go doWork(i, wg)
	}
	wg.Wait()
	fmt.Println("All complete")
}
