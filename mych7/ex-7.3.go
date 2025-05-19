package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func findFactors(number int, wg *sync.WaitGroup) []int {
	result := make([]int, 0)
	for i := 1; i <= number; i++ {
		if number%i == 0 {
			result = append(result, i)
		}
	}
	wg.Done()
	return result
}

func main() {
	resultCh := make(chan []int, 10)
	wg := sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			resultCh <- findFactors(rand.Intn(1000), &wg)
		}()
	}

	wg.Wait()
	close(resultCh)

	for res := range resultCh {
		fmt.Println(res)
	}
}
