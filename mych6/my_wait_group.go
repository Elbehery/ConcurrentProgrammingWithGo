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

type MyWaitGroup struct {
	sema *semaphore
}

func NewMyWaitGroup(size int) *MyWaitGroup {
	return &MyWaitGroup{sema: newSemaphore(1 - size)}
}

func (wg *MyWaitGroup) Wait() {
	wg.sema.acquire()
}

func (wg *MyWaitGroup) Done() {
	wg.sema.release()
}

func doWork(id int, wg *MyWaitGroup) {
	fmt.Println(id, "Done working ")
	wg.Done()
}

func main() {
	wg := NewMyWaitGroup(4)
	for i := 1; i <= 4; i++ {
		go doWork(i, wg)
	}
	wg.Wait()
	fmt.Println("All complete")
}
