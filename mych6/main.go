package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func DoWork(id int, wg *sync.WaitGroup) {
	n := rand.Intn(5)
	time.Sleep(time.Duration(n) * time.Second)
	fmt.Println(id, "Done after", n, " seconds")
	wg.Done()
}

func main() {
	wg := sync.WaitGroup{}
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go DoWork(i, &wg)
	}

	wg.Wait()
	fmt.Println("ALl Done !!!")
}
