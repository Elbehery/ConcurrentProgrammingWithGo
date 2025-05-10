package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

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
	if rwl.readersLock.TryLock() {
		global := true
		if rwl.readerCounter == 0 {
			global = rwl.globalLock.TryLock()
		}
		if global {
			rwl.readerCounter++
		}
		rwl.readersLock.Unlock()
		return global
	} else {
		return false
	}
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

// =====================================

/*
Note: this program has a race condition for demonstration purposes
Additionally we have a timer at the end which you might need to adjust
depending on how fast your internet connection is.
In later chapters we cover how to wait for threads to complete their work
*/
func countLetters(url string, frequency map[string]int, mutex *sync.Mutex) {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		panic("Server's error: " + resp.Status)
	}
	body, _ := io.ReadAll(resp.Body)
	wordRegex := regexp.MustCompile(`[a-zA-Z]+`)
	mutex.Lock()
	for _, word := range wordRegex.FindAllString(string(body), -1) {
		wordLower := strings.ToLower(word)
		frequency[wordLower] += 1
	}
	mutex.Unlock()
	fmt.Println("Completed:", url)
}

func main() {
	mutex := sync.Mutex{}
	var frequency = make(map[string]int)
	for i := 1000; i <= 1020; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		go countLetters(url, frequency, &mutex)
	}
	time.Sleep(10 * time.Second)
	mutex.Lock()
	for k, v := range frequency {
		fmt.Println(k, "->", v)
	}
	mutex.Unlock()
}
