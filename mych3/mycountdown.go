package main

import (
	"fmt"
	"time"
)

func MyCountDown(seconds *int) {
	for *seconds > 0 {
		time.Sleep(1 * time.Second)
		*seconds -= 1
	}
}

func main() {
	counter := 10
	go MyCountDown(&counter)

	for counter > 0 {
		time.Sleep(500 * time.Millisecond)
		fmt.Println(counter)
	}
}
