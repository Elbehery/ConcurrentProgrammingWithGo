package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	pattern := os.Args[1]
	dirPath := os.Args[2]
	wg := &sync.WaitGroup{}

	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if !f.IsDir() {
			wg.Add(1)
			go grepDir(pattern, filepath.Join(dirPath, f.Name()), wg)
		}
	}

	wg.Wait()
}

func grepDir(pattern, fileName string, wg *sync.WaitGroup) {
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
