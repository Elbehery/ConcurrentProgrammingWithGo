package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

func main() {
	targetString := os.Args[1]
	fileNames := os.Args[2:]
	wg := &sync.WaitGroup{}

	for _, f := range fileNames {
		wg.Add(1)
		go grep(targetString, f, wg)
	}

	wg.Wait()
}

func grep(target, fileName string, wg *sync.WaitGroup) {
	content, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	if strings.Index(string(content), target) != -1 {
		fmt.Printf("%s contains a match", fileName)
	}
	wg.Done()
}
