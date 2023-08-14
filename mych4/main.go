package main

import "sync"

type ReaderWriterLock struct {
	readersCounter int
	readersLock    sync.Mutex
	globalLock     sync.Mutex
}

func (l *ReaderWriterLock) RLock() {
	l.readersLock.Lock()
	l.readersCounter++
	if l.readersCounter == 1 {
		l.globalLock.Lock()
	}
	l.readersLock.Unlock()
}

func (l *ReaderWriterLock) RUnLock() {
	l.readersLock.Lock()
	l.readersCounter--
	if l.readersCounter == 0 {
		l.globalLock.Unlock()
	}
	l.readersLock.Unlock()
}

func (l *ReaderWriterLock) Lock() {
	l.globalLock.Lock()
}

func (l *ReaderWriterLock) UnLock() {
	l.globalLock.Unlock()
}
