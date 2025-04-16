package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var x int64 = 10000000
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100000; j++ {
				atomic.AddInt64(&x, -1)
			}
		}()
	}

	wg.Wait()

	fmt.Printf("Result: %d", atomic.LoadInt64(&x))
}
