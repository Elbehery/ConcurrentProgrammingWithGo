package main

import "sync"

type MySemaphore struct {
	permits int
	cond    *sync.Cond
}

func NewMySemaphore(n int) *MySemaphore {
	return &MySemaphore{permits: n, cond: sync.NewCond(&sync.Mutex{})}
}

func (s *MySemaphore) Acquire() {
	s.cond.L.Lock()
	if s.permits <= 0 {
		s.cond.Wait()
	}
	s.permits--
	s.cond.L.Unlock()
}

func (s *MySemaphore) Release() {
	s.cond.L.Lock()
	s.permits++
	s.cond.Signal()
	s.cond.L.Unlock()
}
