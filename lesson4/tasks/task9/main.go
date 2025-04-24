package main

import (
	"context"
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

func (s *Semaphore) Acquire(ctx context.Context) error {
	select {
	case s.ch <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *Semaphore) Release() {
	select {
	case <-s.ch:
	default:

	}
}

func main() {
	s := NewSemaphore(2)
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(i)*300*time.Millisecond)
			defer cancel()

			fmt.Printf("Goroutine %d: waiting for slot\n", i)
			if err := s.Acquire(ctx); err != nil {
				fmt.Printf("Goroutine %d: canceled or timed out: %v\n", i, err)
				return
			}

			defer s.Release()
			fmt.Printf("Goroutine %d: acquired slot\n", i)
			time.Sleep(1 * time.Second)
			fmt.Printf("Goroutine %d: done\n", i)
		}(i)
	}

	wg.Wait()
}
