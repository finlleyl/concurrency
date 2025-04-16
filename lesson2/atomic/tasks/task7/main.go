package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var counter int64
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 10000; j++ {
				atomic.AddInt64(&counter, 1)
				fmt.Printf("Counter: %d\n from goroutine %d\n", counter, i)
			}
		}(i)
	}

	wg.Wait()

	fmt.Printf("%d", atomic.LoadInt64(&counter))
}
