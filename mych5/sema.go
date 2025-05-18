package main

import (
	"fmt"
	"sync"
)

type Semaphore struct {
	permits int
	cnd     *sync.Cond
}

func NewSemaphore(n int) *Semaphore {
	return &Semaphore{
		permits: n,
		cnd:     sync.NewCond(&sync.Mutex{}),
	}
}

func (s *Semaphore) Acquire() {
	s.cnd.L.Lock()
	for s.permits <= 0 {
		s.cnd.Wait()
	}
	s.permits--
	s.cnd.L.Unlock()
}

func (s *Semaphore) Release() {
	s.cnd.L.Lock()
	s.permits++
	s.cnd.Signal()
	s.cnd.L.Unlock()
}

func main() {
	semaphore := NewSemaphore(0)
	for i := 0; i < 50000; i++ {
		go doWork(semaphore)
		fmt.Println("Waiting for child goroutine")
		semaphore.Acquire()
		fmt.Println("Child goroutine finished")
	}
}

func doWork(semaphore *Semaphore) {
	fmt.Println("Work started")
	fmt.Println("Work finished")
	semaphore.Release()
}
