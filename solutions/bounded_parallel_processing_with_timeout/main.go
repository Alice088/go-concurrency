package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type outVal struct {
	val int
	err error
}

var errTimeout = errors.New("timeout")

func processData(ctx context.Context, val int) <-chan outVal {
	ch := make(chan outVal)

	go func(ctx context.Context) {
		time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
		ch <- outVal{val: val * 2}
		close(ch)
	}(ctx)

	return ch
}

func worker(ctx context.Context, wg *sync.WaitGroup, in <-chan int, out chan<- outVal) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			out <- outVal{err: errTimeout} //for timeout demonstration
			return
		case val, ok := <-in:
			if !ok {
				return
			}
			select {
			case <-ctx.Done():
				out <- outVal{err: errTimeout} //for timeout demonstration
				return
			case v := <-processData(ctx, val):
				out <- v
			}
		}
	}
}

func processParallel(ctx context.Context, in <-chan int, out chan<- outVal, numWorkers int) {
	wg := &sync.WaitGroup{}

	wg.Add(numWorkers)
	for range numWorkers {
		go worker(ctx, wg, in, out)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
}

func main() {
	in := make(chan int)
	out := make(chan outVal)

	go func() {
		for i := range 100 {
			in <- i + 1
		}
		close(in)
	}()

	now := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	processParallel(ctx, in, out, 5)

	i := 1
	for val := range out {
		fmt.Printf("g%d) v=%d, err=%v\n", i, val.val, val.err)
		i++
	}

	fmt.Println(time.Since(now))
}
