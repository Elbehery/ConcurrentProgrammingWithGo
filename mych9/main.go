package main

import "fmt"

func printNums(numsCh <-chan int, quitCh chan<- struct{}) {
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(<-numsCh)
		}

		close(quitCh)
	}()
}

func main() {
	numsCh := make(chan int)
	quitCh := make(chan struct{})

	printNums(numsCh, quitCh)
	next := 0

	for i := 1; ; i++ {
		next += i
		select {
		case numsCh <- next:
		case <-quitCh:
			fmt.Println("DONE !!")
			return
		}
	}
}
