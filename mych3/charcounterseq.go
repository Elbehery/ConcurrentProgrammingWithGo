package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

const allLetters = "abcdefghijklmnopqrstuvwxyz"

type ConcurrentFrequencyTable struct {
	table []int
	mu    sync.RWMutex
}

func NewConcurrentFrequencyTable() *ConcurrentFrequencyTable {
	return &ConcurrentFrequencyTable{
		table: make([]int, 26),
		mu:    sync.RWMutex{},
	}
}

func CountChars(url string, freqTable *ConcurrentFrequencyTable, wg *sync.WaitGroup) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	freqTable.mu.Lock()
	for _, b := range data {
		c := strings.ToLower(string(b))
		index := strings.Index(allLetters, c)
		if index >= 0 {
			freqTable.table[index] += 1
		}
	}
	freqTable.mu.Unlock()
	fmt.Printf("processing %s completed\n", url)
	wg.Done()
}

func main() {
	fmt.Println(time.Now())
	freq := NewConcurrentFrequencyTable()
	wg := sync.WaitGroup{}

	for i := 1000; i <= 1200; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		wg.Add(1)
		CountChars(url, freq, &wg)
	}
	wg.Wait()
	for i, c := range allLetters {
		fmt.Printf("%c-%d\n", c, freq.table[i])
	}
	fmt.Println(time.Now())
}
