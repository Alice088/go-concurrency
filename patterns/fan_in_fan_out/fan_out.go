package fan_in_fan_out

import (
	"context"
	"fmt"
	"time"
)

type WorkerConfig struct {
	WorkerID int
	Jobs     <-chan int
	Out      chan<- string
	Ctx      context.Context
}

func Worker(config WorkerConfig) {
	for job := range config.Jobs {
		time.Sleep(500 * time.Millisecond) // hard work
		select {
		case <-config.Ctx.Done():
			return
		case config.Out <- fmt.Sprintf("Worker #%d -- job #%d -> finished", config.WorkerID, job):
		}
	}
}
