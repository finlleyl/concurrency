package main

import (
	"context"
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
