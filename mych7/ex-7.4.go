package main

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

type ChannelCnd[M any] struct {
	cnd      *sync.Cond
	buf      *list.List
	capacity int
}

func NewChannelCnd[M any](capacity int) *ChannelCnd[M] {
	return &ChannelCnd[M]{
		cnd:      sync.NewCond(&sync.Mutex{}),
		buf:      list.New(),
		capacity: capacity,
	}
}

func (c *ChannelCnd[M]) Send(msg M) {
	c.cnd.L.Lock()
	for c.buf.Len() == c.capacity {
		c.cnd.Wait()
	}
	c.buf.PushBack(msg)
	c.cnd.Broadcast()
	c.cnd.L.Unlock()
}

func (c *ChannelCnd[M]) Receive() M {
	c.cnd.L.Lock()
	for c.buf.Len() == 0 {
		c.cnd.Wait()
	}
	v := c.buf.Remove(c.buf.Front()).(M)
	c.cnd.Broadcast()
	c.cnd.L.Unlock()

	return v
}

func receiverChanCnd(messages *ChannelCnd[int], wGroup *sync.WaitGroup) {
	msg := 0
	for msg != -1 {
		time.Sleep(1 * time.Second)
		msg = messages.Receive()
		fmt.Println("Received:", msg)
	}
	wGroup.Done()
}

func main() {
	channel := NewChannelCnd[int](2)
	wGroup := sync.WaitGroup{}
	wGroup.Add(1)
	go receiverChanCnd(channel, &wGroup)
	for i := 1; i <= 6; i++ {
		fmt.Println("Sending: ", i)
		channel.Send(i)
	}
	channel.Send(-1)
	wGroup.Wait()
}
