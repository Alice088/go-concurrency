package main

import (
	"concurrency/common"
	"concurrency/patterns/fan_in_fan_out"
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	channels := common.GenerateWorkersWithDuration(ctx, 3, time.Second)

	out := fan_in_fan_out.FanIN(ctx, channels)

	for v := range out {
		fmt.Println(v)
	}
}
