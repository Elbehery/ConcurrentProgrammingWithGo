package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	fileNames := os.Args[1:]
	for _, f := range fileNames {
		go printContent(f)
	}

	time.Sleep(5 * time.Second)
}

func printContent(fileName string) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))
}
