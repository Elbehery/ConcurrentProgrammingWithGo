package main

import (
	"fmt"
	"sync"
	"time"
)

type MyBarrier struct {
	size, waitCount int
	cnd             *sync.Cond
}

func NewMyBarrier(size int) *MyBarrier {
	return &MyBarrier{
		size: size,
		cnd:  sync.NewCond(&sync.Mutex{}),
	}
}

func (b *MyBarrier) Wait() {
	b.cnd.L.Lock()
	b.waitCount++
	if b.waitCount == b.size {
		b.waitCount = 0
		b.cnd.Broadcast()
	} else {
		b.cnd.Wait()
	}
	b.cnd.L.Unlock()
}

func workAndWait(name string, timeToWork int, barrier *MyBarrier) {
	start := time.Now()
	for {
		fmt.Println(time.Since(start), name, "is running")
		time.Sleep(time.Duration(timeToWork) * time.Second)
		fmt.Println(time.Since(start), name, "is waiting on barrier")
		barrier.Wait()
	}
}

func main() {
	barrier := NewMyBarrier(2)
	go workAndWait("Red", 2, barrier)
	go workAndWait("Blue", 1, barrier)
	time.Sleep(10 * time.Second)
}
