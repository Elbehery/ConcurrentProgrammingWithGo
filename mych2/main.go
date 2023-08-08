package main

import (
	"fmt"
	"sync"
	"time"
)

func doWork(id int, wg *sync.WaitGroup) {
	fmt.Printf("Work %d started at %s\n", id, time.Now().Format("15:04:05"))
	time.Sleep(1 * time.Second)
	fmt.Printf("Work %d finished at %s\n", id, time.Now().Format("15:04:05"))
	wg.Done()
}

func main() {

	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go doWork(i, &wg)
	}
	wg.Wait()
}
