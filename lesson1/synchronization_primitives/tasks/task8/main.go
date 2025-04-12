package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var mu sync.Mutex
	cond := sync.NewCond(&mu)

	var wg sync.WaitGroup
	flag := false

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			mu.Lock()
			for !flag {
				cond.Wait()
			}
			fmt.Printf("Горутина %d стартует\n", id)
			mu.Unlock()
		}(i)
	}

	time.Sleep(2 * time.Second)

	mu.Lock()
	flag = true
	cond.Broadcast() // разбудить всех, кто ждал
	mu.Unlock()

	fmt.Println("main: подал сигнал Broadcast()")

	wg.Wait()
}
