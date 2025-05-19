package main

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

type semaphore struct {
	permits int
	cnd     *sync.Cond
}

func newSemaphore(permits int) *semaphore {
	return &semaphore{
		permits: permits,
		cnd:     sync.NewCond(&sync.Mutex{}),
	}
}

func (s *semaphore) acquire() {
	s.cnd.L.Lock()
	for s.permits <= 0 {
		s.cnd.Wait()
	}
	s.permits--
	s.cnd.L.Unlock()
}

func (s *semaphore) release() {
	s.cnd.L.Lock()
	s.permits++
	s.cnd.Signal()
	s.cnd.L.Unlock()
}

type MyChannel[M any] struct {
	buf               *list.List
	capSema, sizeSema *semaphore
	mux               sync.Mutex
}

func NewMyChannel[M any](capacity int) *MyChannel[M] {
	return &MyChannel[M]{
		buf:      list.New(),
		capSema:  newSemaphore(capacity),
		sizeSema: newSemaphore(0),
		mux:      sync.Mutex{},
	}
}

func (c *MyChannel[M]) Send(m any) {
	c.capSema.acquire()

	c.mux.Lock()
	c.buf.PushBack(m)
	c.mux.Unlock()

	c.sizeSema.release()
}

func (c *MyChannel[M]) Receive() M {
	c.sizeSema.acquire()

	c.mux.Lock()
	v := c.buf.Remove(c.buf.Front()).(M)
	c.mux.Unlock()

	c.capSema.release()

	return v
}

func receiver(messages *MyChannel[int], wGroup *sync.WaitGroup) {
	msg := 0
	for msg != -1 {
		time.Sleep(1 * time.Second)
		msg = messages.Receive()
		fmt.Println("Received:", msg)
	}
	wGroup.Done()
}

func main() {
	channel := NewMyChannel[int](10)
	wGroup := sync.WaitGroup{}
	wGroup.Add(1)
	go receiver(channel, &wGroup)
	for i := 1; i <= 6; i++ {
		fmt.Println("Sending: ", i)
		channel.Send(i)
	}
	channel.Send(-1)
	wGroup.Wait()
}
