/*
	we have var called nums = []int{12,54,89,21,66,47,14,285,96,...}
	we want to calculate the sum of each number by the power of 2 (12^2+54^2+..)
	we want a pipeline to send this numbers to a channel and have 2 concurrent functions
	(worker) receives from the pipeline and do the calc
*/
package main

import (
	"fmt"
	"sync"
)

var (
	mu  = sync.Mutex{}
	wg  = sync.WaitGroup{}
	sum = 0
)

func main() {
	nums := []int{12, 54, 89, 21, 66, 47, 14, 285, 96, 0}
	numChan := make(chan int, len(nums))

	ss := 0
	for _, val := range nums {
		ss += val * val
	}
	fmt.Println(ss)

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go calcWorker(numChan)
	}

	for _, val := range nums {
		numChan <- val
	}
	close(numChan)

	wg.Wait()

	fmt.Println(sum)
}

func calcWorker(ch <-chan int) {
	defer wg.Done()

	for num := range ch {
		mu.Lock()
		sum += num * num
		mu.Unlock()
	}
}
