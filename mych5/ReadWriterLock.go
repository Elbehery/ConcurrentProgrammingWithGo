package main

import (
	"sync"
)

type ReadWriteLock struct {
	activeReaders  int
	waitingWriters int
	activeWriter   bool
	cond           *sync.Cond
}

func NewReadWriteLock() *ReadWriteLock {
	return &ReadWriteLock{cond: sync.NewCond(&sync.Mutex{})}
}

func (rwl *ReadWriteLock) ReadLock() {
	rwl.cond.L.Lock()
	for rwl.waitingWriters > 0 || rwl.activeWriter {
		rwl.cond.Wait()
	}
	rwl.activeReaders++
	rwl.cond.L.Unlock()
}

func (rwl *ReadWriteLock) WriteLock() {
	rwl.cond.L.Lock()
	rwl.waitingWriters++
	for rwl.activeWriter || rwl.activeReaders > 0 {
		rwl.cond.Wait()
	}
	rwl.waitingWriters--
	rwl.activeWriter = true
	rwl.cond.L.Unlock()
}

func (rwl *ReadWriteLock) ReadUnLock() {
	rwl.cond.L.Lock()
	rwl.activeReaders--
	if rwl.activeReaders == 0 {
		rwl.cond.Broadcast()
	}
	rwl.cond.L.Unlock()
}

func (rwl *ReadWriteLock) WriteUnLock() {
	rwl.cond.L.Lock()
	rwl.activeWriter = false
	rwl.cond.Broadcast()
	rwl.cond.L.Unlock()
}

//func main() {
//	rwMutex := NewReadWriteLock()
//	for i := 0; i < 2; i++ {
//		go func() {
//			for {
//				rwMutex.ReadLock()
//				time.Sleep(1 * time.Second)
//				fmt.Println("Read done")
//				rwMutex.ReadUnLock()
//			}
//		}()
//	}
//	time.Sleep(1 * time.Second)
//	rwMutex.WriteLock()
//	fmt.Println("Write finished")
//}
