package fan_in_fan_out

import (
	"context"
	"sync"
)

func FanIN(ctx context.Context, channels []<-chan int) <-chan int {
	out := make(chan int)

	go func() {
		wg := &sync.WaitGroup{}

		wg.Add(len(channels))
		for _, ch := range channels {
			go func() {
				defer wg.Done()

				for {
					select {
					case <-ctx.Done():
						return
					case v, ok := <-ch:
						if !ok {
							return
						}
						select {
						case <-ctx.Done():
							return
						case out <- v:
						}
					}
				}
			}()
		}

		wg.Wait()
		close(out)
	}()

	return out
}
