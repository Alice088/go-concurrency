package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	t := time.Now()

	nums := make(chan int)
	switchCh := make(chan bool)
	var slice []int

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	go closeCh(ctx, nums)

	go func() {
		switchCh <- true
	}()

	go wEven(nums, switchCh)
	go wOdd(nums, switchCh)

	for num := range nums {
		slice = append(slice, num)
	}

	fmt.Println(slice)
	fmt.Println(time.Since(t).String())
}

func wEven(nums chan int, switchCh chan bool) {
	for i := range 10 {
		select {
		case v := <-switchCh:
			if v == true {
				if i%2 == 0 {
					fmt.Printf("wEven: write: %d\n", i+1)
					nums <- i + 1
					switchCh <- false
					continue
				}
			}
			switchCh <- v
		}
	}
}

func wOdd(nums chan int, switchCh chan bool) {
	for i := range 10 {
		select {
		case v := <-switchCh:
			if v == false {
				if i%2 != 0 {
					fmt.Printf("wOdd: write: %d\n", i+1)
					nums <- i + 1
					switchCh <- true
					continue
				}
			}
			switchCh <- v
		}
	}
	close(nums)
}

func closeCh(ctx context.Context, ch chan int) {
	select {
	case <-ctx.Done():
		close(ch)
	}
}
