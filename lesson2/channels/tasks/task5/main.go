package main

import (
	"fmt"
	"sync"
	"time"
)

var done = make(chan struct{})

func worker(wg *sync.WaitGroup, i int) {
	for {
		select {
		case <-done:
			fmt.Println("Завершаем", i)
			wg.Done()
			return
		default:
			fmt.Println(i)
		}
	}
}

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go worker(&wg, i)
	}

	time.Sleep(1 * time.Second)
	close(done)
	wg.Wait()
}
