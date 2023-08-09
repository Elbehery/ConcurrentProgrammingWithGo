package main

import (
	"fmt"
	"log"
	"os"
	"sync"
)

func PrintFileStdOut(name string, wg *sync.WaitGroup) {
	f, err := os.ReadFile(name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(f))
	wg.Done()
}

func main() {
	files := os.Args[1:]
	wg := sync.WaitGroup{}
	for _, f := range files {
		wg.Add(1)
		go PrintFileStdOut(f, &wg)
	}
	wg.Wait()
}
