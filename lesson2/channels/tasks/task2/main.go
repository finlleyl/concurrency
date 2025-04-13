package main

import (
	"fmt"
	"sync"
)

func main() {
	ch1 := make(chan int)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			ch1 <- i
			wg.Done()
		}(i)
	}

	go func() {
		wg.Wait()
		close(ch1)
	}()

	for v := range ch1 {
		fmt.Println(v)
	}

	v, ok := <-ch1
	fmt.Println(v, ok)
}
