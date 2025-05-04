package main

import (
	"fmt"
	"log"
	"os"
	"sync"
)

func main() {
	fileNames := os.Args[1:]
	var wg sync.WaitGroup

	for _, f := range fileNames {
		wg.Add(1)
		go printContent(f, &wg)
	}

	wg.Wait()
}

func printContent(fileName string, wg *sync.WaitGroup) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))
	wg.Done()
}
