package main

import (
	"concurrency/patterns/fan_in_fan_out"
	"fmt"
	"sync"
)

func main() {
	jobs := make(chan int)
	out := make(chan string)
	wg := &sync.WaitGroup{}
	go func() {
		for i := range 20 {
			jobs <- i + 1
		}
		close(jobs)
	}()

	wg.Add(4)
	for i := range 4 {
		go func() {
			defer wg.Done()
			fan_in_fan_out.Worker(i+1, jobs, out)
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	for v := range out {
		fmt.Println(v)
	}
}
