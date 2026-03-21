package main

import (
	"concurrency/patterns/fan_in_fan_out"
	"context"
	"fmt"
	"sync"
	"time"
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	wg.Add(4)
	for i := range 4 {
		go func() {
			defer wg.Done()
			fan_in_fan_out.Worker(fan_in_fan_out.WorkerConfig{
				WorkerID: i + 1,
				Jobs:     jobs,
				Ctx:      ctx,
				Out:      out,
			})
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
