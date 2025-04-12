package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()

			t := rand.Intn(100)
			sleepTime := time.Duration(t) * time.Millisecond
			time.Sleep(sleepTime)

			fmt.Printf("Hello from goroutine %d\n", n)
		}(i)
	}

	wg.Wait()
	fmt.Println("All goroutines finished!")

}
