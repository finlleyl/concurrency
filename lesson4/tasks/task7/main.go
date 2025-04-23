package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Stage func(context.Context, <-chan int) <-chan int

func generator(ctx context.Context, nums ...int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for _, n := range nums {
			select {
			case out <- n:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

func heavy(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for n := range in {
			select {
			case out <- n * n:
				time.Sleep(200 * time.Millisecond)
			case <-ctx.Done():
				return
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
	out := make(chan int)
	for _, ch := range chans {
		wg.Add(1)
		ch := ch
		go func(ch <-chan int) {
			defer wg.Done()
			for n := range ch {
				select {
				case out <- n:
				case <-ctx.Done():
					return
				}
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func pipeline(ctx context.Context, in <-chan int, stages ...Stage) <-chan int {
	for _, stage := range stages {
		in = stage(ctx, in)
	}

	return in
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	gen := generator(ctx, nums...)

	fanned := func(ctx context.Context, in <-chan int) <-chan int {
		outChans := fanOut(ctx, in, 3)
		return fanIn(ctx, outChans...)
	}

	out := pipeline(ctx,
		gen, fanned)

	for n := range out {
		fmt.Println(n)
	}
}
