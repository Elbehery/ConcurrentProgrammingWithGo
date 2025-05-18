package main

import (
	"fmt"
	"sync"
	"time"
)

type ReadWriteMutex struct {
	currentReaders, waitingWriters int
	activeWriter                   bool
	cond                           *sync.Cond
}

func NewReadWriteMutex() *ReadWriteMutex {
	return &ReadWriteMutex{
		cond: sync.NewCond(&sync.Mutex{}),
	}
}

func (m *ReadWriteMutex) Lock() {
	m.cond.L.Lock()
	m.waitingWriters++
	for m.currentReaders > 0 || m.activeWriter {
		m.cond.Wait()
	}
	m.waitingWriters--
	m.activeWriter = true
	m.cond.L.Unlock()
}

func (m *ReadWriteMutex) UnLock() {
	m.cond.L.Lock()
	m.activeWriter = false
	m.cond.Broadcast()
	m.cond.L.Unlock()
}

func (m *ReadWriteMutex) RLock() {
	m.cond.L.Lock()
	for m.activeWriter || m.waitingWriters > 0 {
		m.cond.Wait()
	}
	m.currentReaders++
	m.cond.L.Unlock()
}

func (m *ReadWriteMutex) RUnLock() {
	m.cond.L.Lock()
	m.currentReaders--
	if m.currentReaders == 0 {
		m.cond.Broadcast()
	}
	m.cond.L.Unlock()
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
