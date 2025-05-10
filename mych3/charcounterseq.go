package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func countLetters(url string, freq []int, wg *sync.WaitGroup) {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	defer wg.Done()

	if resp.StatusCode != 200 {
		log.Fatal(resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	for _, b := range data {
		c := strings.ToLower(string(b))
		idx := strings.Index(alphabet, c)
		if idx >= 0 {
			freq[idx]++
		}
	}
}

func main() {
	wg := &sync.WaitGroup{}
	freq := make([]int, 26)
	for i := 1000; i <= 1030; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		wg.Add(1)
		go countLetters(url, freq, wg)
	}

	wg.Wait()

	for i, c := range alphabet {
		fmt.Printf("%c-%d\n", c, freq[i])
	}
}
