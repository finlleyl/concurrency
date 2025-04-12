package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var mu sync.Mutex
	cond := sync.NewCond(&mu)

	room := false

	enterRoom := func(id int, wg *sync.WaitGroup) {
		defer wg.Done()

		mu.Lock()
		for room {
			cond.Wait()
		}

		room = true
		fmt.Printf("Горутина %d вошла в комнату\n", id)
		mu.Unlock()

		time.Sleep(200 * time.Millisecond)

		mu.Lock()

		room = false
		fmt.Printf("Горутина %d вышла из комнаты\n", id)
		cond.Signal()
		mu.Unlock()
	}

	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go enterRoom(i, &wg)
	}

	wg.Wait()
}
