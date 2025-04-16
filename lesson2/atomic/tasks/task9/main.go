package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var flag atomic.Bool
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			if !flag.Load() {
				flag.Store(true)
				fmt.Printf("Goroutine changes flag to true\n")
			} else {
				flag.Store(false)
				fmt.Printf("Goroutine changes flag to false\n")
			}
		}()
	}
}
