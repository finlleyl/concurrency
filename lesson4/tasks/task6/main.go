package main

import (
	"context"
	"sync"
	"time"
)

func heavy(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				time.Sleep(500 * time.Millisecond)
				select {
				case <-ctx.Done():
					return
				case out <- v * 2:
				}
			}
		}
	}()
	return out
}

func fanOut(ctx context.Context, in <-chan int, workers int) []<-chan int {
	out := make([]<-chan int, workers)

	for i := 0; i < workers; i++ {
		out[i] = heavy(ctx, in)
	}

	return out
}

func fanIn(ctx context.Context, chans ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	finalCh := make(chan int)

	wg.Add(len(chans))
	for _, c := range chans {
		c := c
		go func(c <-chan int) {
			defer wg.Done()
			for v := range c {
				select {
				case <-ctx.Done():
					return
				case finalCh <- v:
				}
			}
		}(c)
	}

	go func() {
		wg.Wait()
		close(finalCh)
	}()

	return finalCh
}
