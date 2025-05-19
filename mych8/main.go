package main

import (
	"fmt"
	"time"
)

func writeEvery(msg string, duration time.Duration) <-chan string {
	ch := make(chan string)
	go func() {
		for {
			time.Sleep(duration)
			ch <- msg
		}
	}()

	return ch
}

func main() {
	chA := writeEvery("A", 1*time.Second)
	chB := writeEvery("B", 2*time.Second)

	for {
		select {
		case a := <-chA:
			fmt.Println(a)
		case b := <-chB:
			fmt.Println(b)
		}
	}
}
