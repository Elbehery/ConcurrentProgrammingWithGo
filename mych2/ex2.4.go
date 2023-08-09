package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func MyGrepDir(path, pattern string, file os.DirEntry, wg *sync.WaitGroup) {
	fullPath := filepath.Join(path, file.Name())
	if file.IsDir() {
		files, err := os.ReadDir(fullPath)
		if err != nil {
			log.Fatal(err)
		}
		for _, f := range files {
			wg.Add(1)
			go GrepDir(path, pattern, f, wg)
		}
	} else {
		f, err := os.ReadFile(fullPath)
		if err != nil {
			log.Fatal(err)
		}
		if strings.Contains(string(f), pattern) {
			fmt.Println(fullPath, "contains a match with", pattern)
		} else {
			fmt.Println(fullPath, "does NOT contain a match with", pattern)
		}
	}
	wg.Done()
}

func main() {
	pat := os.Args[1]
	path := os.Args[2]

	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	wg := sync.WaitGroup{}
	for _, f := range files {
		wg.Add(1)
		go MyGrepDir(path, pat, f, &wg)
	}
	wg.Wait()
}
