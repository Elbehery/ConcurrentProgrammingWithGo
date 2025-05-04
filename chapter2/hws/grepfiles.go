package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

func main() {
	pattern := os.Args[1]
	fileNames := os.Args[2:]
	wg := &sync.WaitGroup{}

	for _, f := range fileNames {
		wg.Add(1)
		go grep(pattern, f, wg)
	}

	wg.Wait()
}

func grep(pattern, fileName string, wg *sync.WaitGroup) {
	content, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	if strings.Contains(string(content), pattern) {
		fmt.Printf("%s contains a match with %s \n", fileName, pattern)
	} else {
		fmt.Printf("%s does not contains a match with %s \n", fileName, pattern)
	}
	wg.Done()
}
