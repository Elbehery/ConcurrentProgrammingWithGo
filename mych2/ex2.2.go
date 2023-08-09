package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

func MyGrep(pattern string, filename string, wg *sync.WaitGroup) {
	f, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	if strings.Contains(string(f), pattern) {
		fmt.Println(filename, "contains a match with", pattern)
	} else {
		fmt.Println(filename, "does NOT contain a match with", pattern)
	}
	wg.Done()
}

func main() {
	pattern := os.Args[1]
	files := os.Args[2:]
	wg := sync.WaitGroup{}
	for _, f := range files {
		wg.Add(1)
		go MyGrep(pattern, f, &wg)
	}
	wg.Wait()
}
