package main

import (
	"fmt"
	"sync"
	"time"
)

type MyWSemaphore struct {
	permits int
	cond    *sync.Cond
}

func NewMyWSemaphore(n int) *MyWSemaphore {
	return &MyWSemaphore{
		permits: n,
		cond:    sync.NewCond(&sync.Mutex{}),
	}
}

func (s *MyWSemaphore) Acquire(n int) {
	s.cond.L.Lock()
	if s.permits-n < 0 {
		s.cond.Wait()
	}
	s.permits -= n
	s.cond.L.Unlock()
}

func (s *MyWSemaphore) Release(n int) {
	s.cond.L.Lock()
	s.permits += n
	s.cond.Signal()
	s.cond.L.Unlock()
}

func main() {
	sema := NewMyWSemaphore(3)
	sema.Acquire(2)
	fmt.Println("Parent thread acquired semaphore")
	go func() {
		sema.Acquire(2)
		fmt.Println("Child thread acquired semaphore")
		sema.Release(2)
		fmt.Println("Child thread released semaphore")
	}()
	time.Sleep(3 * time.Second)
	fmt.Println("Parent thread releasing semaphore")
	sema.Release(2)
	time.Sleep(1 * time.Second)
}
