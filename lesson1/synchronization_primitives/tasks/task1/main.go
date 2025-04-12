package main

import "sync"

func main() {
	c := 0

	var wg sync.WaitGroup
	var mx sync.Mutex

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mx.Lock()
			c++
			mx.Unlock()
		}()
	}

	wg.Wait()
}
