package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	var mx sync.RWMutex
	var wg sync.WaitGroup
	var data int
	go func() {
		for i := 0; i < 100; i++ {
			mx.Lock()
			data = rand.Int()
			mx.Unlock()
			time.Sleep(300 * time.Millisecond)
		}
	}()

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mx.RLock()
			fmt.Println(data)
			mx.RUnlock()
		}()
	}

	wg.Wait()
}
