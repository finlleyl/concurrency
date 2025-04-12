package main

import (
	"math/rand"
	"sync"
)

func main() {
	var mx sync.Mutex
	var wg sync.WaitGroup

	res := []int{}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mx.Lock()
			res = append(res, rand.Int())
			mx.Unlock()
		}()
	}

	wg.Wait()
	
}
