package main

import "sync"

type ReaderWriterLock struct {
	readerCounter           int
	readersLock, globalLock sync.Mutex
}

func NewReaderWriterLock() *ReaderWriterLock {
	return &ReaderWriterLock{
		readerCounter: 0,
		readersLock:   sync.Mutex{},
		globalLock:    sync.Mutex{},
	}
}

func (rwl *ReaderWriterLock) Lock() {
	rwl.globalLock.Lock()
}

func (rwl *ReaderWriterLock) TryLock() bool {
	return rwl.globalLock.TryLock()
}

func (rwl *ReaderWriterLock) RLock() {
	rwl.readersLock.Lock()
	rwl.readerCounter++
	if rwl.readerCounter == 1 {
		rwl.globalLock.Lock()
	}
	rwl.readersLock.Unlock()
}

func (rwl *ReaderWriterLock) TryReadLock() bool {
	return rwl.readersLock.TryLock()
}

func (rwl *ReaderWriterLock) UnLock() {
	rwl.globalLock.Unlock()
}

func (rwl *ReaderWriterLock) RUnLock() {
	rwl.readersLock.Lock()
	rwl.readerCounter--
	if rwl.readerCounter == 0 {
		rwl.globalLock.Unlock()
	}
	rwl.readersLock.Unlock()
}
