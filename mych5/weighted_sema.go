package main

import (
	"fmt"
	"sync"
	"time"
)

type WeightedSemaphore struct {
	permits int
	cnd     *sync.Cond
}

func NewWeightedSemaphore(n int) *WeightedSemaphore {
	return &WeightedSemaphore{
		permits: n,
		cnd:     sync.NewCond(&sync.Mutex{}),
	}
}

func (w *WeightedSemaphore) Acquire(permits int) {
	w.cnd.L.Lock()
	for w.permits-permits <= 0 {
		w.cnd.Wait()
	}
	w.permits -= permits
	w.cnd.L.Unlock()
}
func (w *WeightedSemaphore) Release(permits int) {
	w.cnd.L.Lock()
	w.permits += permits
	w.cnd.Broadcast()
	w.cnd.L.Unlock()
}

func main() {
	sema := NewWeightedSemaphore(3)
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
