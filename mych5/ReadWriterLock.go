package main

import (
	"fmt"
	"sync"
	"time"
)

type ReadWriteMutex struct {
	readerCounter, writersWaiting int
	activeWriter                  bool
	cond                          *sync.Cond
}

func NewReadWriteMutex() *ReadWriteMutex {
	return &ReadWriteMutex{
		cond: &sync.Cond{
			L: &sync.Mutex{},
		},
	}
}

func (rwm *ReadWriteMutex) RLock() {
	rwm.cond.L.Lock()
	if rwm.activeWriter || rwm.writersWaiting > 0 {
		rwm.cond.Wait()
	}
	rwm.readerCounter++
	rwm.cond.L.Unlock()
}

func (rwm *ReadWriteMutex) RUnLock() {
	rwm.cond.L.Lock()
	rwm.readerCounter--
	if rwm.readerCounter == 0 {
		rwm.cond.Broadcast()
	}
	rwm.cond.L.Unlock()
}

func (rwm *ReadWriteMutex) Lock() {
	rwm.cond.L.Lock()
	rwm.writersWaiting++
	if rwm.activeWriter || rwm.readerCounter > 0 {
		rwm.cond.Wait()
	}
	rwm.writersWaiting--
	rwm.activeWriter = true
	rwm.cond.L.Unlock()
}

func (rwm *ReadWriteMutex) UnLock() {
	rwm.cond.L.Lock()
	rwm.activeWriter = false
	rwm.cond.Broadcast()
	rwm.cond.L.Unlock()
}

func main() {
	rwMutex := NewReadWriteMutex()
	for i := 0; i < 2; i++ {
		go func() {
			for {
				rwMutex.RLock()
				time.Sleep(1 * time.Second)
				fmt.Println("Read done")
				rwMutex.RUnLock()
			}
		}()
	}
	time.Sleep(1 * time.Second)
	rwMutex.Lock()
	fmt.Println("Write finished")
}
