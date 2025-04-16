package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

var tick atomic.Bool

func main() {
	go func() {
		for !tick.Load() {
			fmt.Printf("tick\n")
			time.Sleep(time.Second)
		}
	}()

	time.Sleep(5 * time.Second)
	tick.Store(true)
}
