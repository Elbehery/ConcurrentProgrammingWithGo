package main

import (
	"fmt"
	"sync"
)

type MyWaitGroup struct {
	size int
	cnd  *sync.Cond
}

func NewMyWaitGroup() *MyWaitGroup {
	return &MyWaitGroup{
		cnd: sync.NewCond(&sync.Mutex{}),
	}
}

func (wg *MyWaitGroup) Add(delta int) {
	wg.cnd.L.Lock()
	wg.size += delta
	wg.cnd.L.Unlock()
}

func (wg *MyWaitGroup) Wait() {
	wg.cnd.L.Lock()
	for wg.size > 0 {
		wg.cnd.Wait()
	}
	wg.cnd.L.Unlock()
}

func (wg *MyWaitGroup) Done() {
	wg.cnd.L.Lock()
	wg.size--
	if wg.size == 0 {
		wg.cnd.Broadcast()
	}
	wg.cnd.L.Unlock()
}

func goWork(id int, wg *MyWaitGroup) {
	fmt.Println(id, "Done working ")
	wg.Done()
}

func main() {
	wg := NewMyWaitGroup()
	for i := 1; i <= 4; i++ {
		wg.Add(2)
		go goWork(i, wg)
		go goWork(i, wg)
	}
	wg.Wait()
	fmt.Println("All complete")
}
