package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func doWork(id int, wg *sync.WaitGroup) {
	fmt.Printf("Work %d started at %s\n", id, time.Now().Format("15:04:05"))
	time.Sleep(1 * time.Second)
	fmt.Printf("Work %d finished at %s\n", id, time.Now().Format("15:04:05"))
	wg.Done()
}

func sayHello() {
	fmt.Println("Hello")
}

func main() {
	go sayHello()
	runtime.Gosched()
	fmt.Println("finished main")
}
