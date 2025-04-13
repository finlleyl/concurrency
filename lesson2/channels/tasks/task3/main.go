package main

import (
	"fmt"
	"sync"
)

func main() {
	ch1 := make(chan int)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		for v := range ch1 {
			fmt.Println(v)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {

		for i := 0; i < 10; i++ {
			ch1 <- i
		}
		close(ch1)
		wg.Done()
	}()

	wg.Wait()

}
