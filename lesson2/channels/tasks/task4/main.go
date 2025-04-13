package main

import (
	"fmt"
	"sync"
)

func main() {
	ch := make(chan string)
	var wg sync.WaitGroup

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ch <- fmt.Sprintf("Hello %d", i)
		}(i)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for v := range ch {
		fmt.Println(v)
	}

}
