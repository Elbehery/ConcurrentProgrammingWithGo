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

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range entries {
		wg.Add(1)
		go grepPath(pattern, dirPath, entry, wg)
	}

	wg.Wait()
}

func grepPath(pattern, path string, entry os.DirEntry, wg *sync.WaitGroup) {
	if entry.IsDir() {
		entries, err := os.ReadDir(entry.Name())
		if err != nil {
			log.Fatal(err)
		}

		for _, e := range entries {
			go grepPath(path, path, e, wg)
		}
	} else {
		content, err := os.ReadFile(filepath.Join(path, entry.Name()))
		if err != nil {
			log.Fatal(err)
		}

		if strings.Contains(string(content), pattern) {
			fmt.Printf("%s contains a match with %s \n", entry.Name(), pattern)
		} else {
			fmt.Printf("%s does not contains a match with %s \n", entry.Name(), pattern)
		}
	}
	wg.Done()
}
