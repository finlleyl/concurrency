package main

import (
	"context"
	"fmt"
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

func square(ctx context.Context, in <-chan int) <-chan int {
	res := make(chan int)
	go func() {
		defer close(res)
		for n := range in {
			select {
			case res <- n * n:
			case <-ctx.Done():
				return
			}
		}
	}()
	return res
}

func increment(ctx context.Context, in <-chan int) <-chan int {
	res := make(chan int)
	go func() {
		defer close(res)
		for n := range in {
			select {
			case res <- n + 1:
			case <-ctx.Done():
				return
			}
		}
	}()
	return res
}

func pipeline(ctx context.Context, in <-chan int, stages ...Stage) <-chan int {
	ch := in
	for _, stage := range stages {
		ch = stage(ctx, ch)
	}
	return ch
}

func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	out := pipeline(ctx,
		generator(ctx, nums...),
		increment,
		square,
	)

	for n := range out {
		fmt.Println(n)
	}
}
