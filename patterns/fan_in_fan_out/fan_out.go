package fan_in_fan_out

import (
	"fmt"
	"time"
)

func Worker(workerID int, jobs <-chan int, out chan<- string) {
	for job := range jobs {
		time.Sleep(500 * time.Millisecond)
		out <- fmt.Sprintf("Worker #%d -- job #%d -> finished", workerID, job)
	}
}
