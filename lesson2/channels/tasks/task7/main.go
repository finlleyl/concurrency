package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func worker(wg *sync.WaitGroup, out chan<- int) {
	defer wg.Done()
	for i := 0; i < 10; i++ {
		out <- rand.Intn(100)
		time.Sleep(200 * time.Millisecond)
	}
}

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	var wg sync.WaitGroup
	wg.Add(2)

	go worker(&wg, ch1)
	go worker(&wg, ch2)

	// отдельная горутина для чтения
	go func() {
		for {
			select {
			case v := <-ch1:
				fmt.Println("ch1:", v)
			case v := <-ch2:
				fmt.Println("ch2:", v)
			}
		}
	}()

	wg.Wait()
}
