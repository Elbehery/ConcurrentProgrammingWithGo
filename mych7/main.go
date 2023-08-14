package main

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

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

type MyChannel[M any] struct {
	capacitySema *sema
	sizeSema     *sema
	mux          sync.Mutex
	buffer       *list.List
}

func NewMyChannel[M any](capacity int) *MyChannel[M] {
	return &MyChannel[M]{
		capacitySema: newSema(capacity),
		sizeSema:     newSema(0),
		buffer:       list.New(),
	}
}

func (c *MyChannel[M]) Send(msg M) {
	c.capacitySema.acquire()
	c.mux.Lock()
	c.buffer.PushBack(msg)
	c.mux.Unlock()
	c.sizeSema.release()
}

func (c *MyChannel[M]) Receive() M {
	c.capacitySema.release()
	c.sizeSema.acquire()
	c.mux.Lock()
	v := c.buffer.Remove(c.buffer.Front()).(M)
	c.mux.Unlock()

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
