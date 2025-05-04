package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func countLetters(url string, freq []int) {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
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
	freq := make([]int, 26)
	for i := 1000; i <= 1030; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		countLetters(url, freq)
	}

	for i, c := range alphabet {
		fmt.Printf("%c-%d\n", c, freq[i])
	}
}
