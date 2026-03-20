package common

import (
	"context"
	"time"
)

func GenerateWorkersWithDuration(ctx context.Context, n int, duration time.Duration) []<-chan int {
	channels := make([]<-chan int, 0, n)

	for i := 0; i < n; i++ {
		ch := make(chan int)

		go func(id int, out chan<- int) {
			defer close(out)

			for j := 0; j < 5; j++ {
				select {
				case <-ctx.Done():
					return
				case out <- id*10 + j:
				}
				time.Sleep(duration)
			}
		}(i, ch)

		channels = append(channels, ch)
	}

	return channels
}
