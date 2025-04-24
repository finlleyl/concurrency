package main

import (
	"fmt"
	"sync"
	"time"
)

type Semaphore struct {
	ch chan struct{}
}

func NewSemaphore(max int) *Semaphore {
	return &Semaphore{make(chan struct{}, max)}
}

func (s *Semaphore) Acquire() {
	s.ch <- struct{}{} // блокирующая операция
}

func (s *Semaphore) Release() {
	<-s.ch
}

func main() {
	s := NewSemaphore(2)
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			fmt.Printf("Goroutine %d waiting for slot...\n", i)
			s.Acquire()
			fmt.Printf("Goroutine %d acquired slot!\n", i)

			time.Sleep(1 * time.Second)

			s.Release()
			fmt.Printf("Goroutine %d released slot.\n", i)
		}(i)
	}

	wg.Wait()
}
