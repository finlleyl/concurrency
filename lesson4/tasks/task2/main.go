package main

import "fmt"

type Stage func(<-chan int) <-chan int

func generator(nums ...int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for _, n := range nums {
			out <- n
		}
	}()

	return out
}

func square(in <-chan int) <-chan int {
	res := make(chan int)

	go func() {
		defer close(res)
		for n := range in {
			res <- n * n
		}
	}()

	return res
}

func increment(in <-chan int) <-chan int {
	res := make(chan int)

	go func() {
		defer close(res)

		for n := range in {
			res <- n + 1
		}
	}()

	return res
}

func pipeline(in <-chan int, stages ...Stage) <-chan int {
	ch := in

	for _, stage := range stages {
		ch = stage(ch)
	}

	return ch
}

func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	out := pipeline(generator(nums...), increment, square)

	for n := range out {
		fmt.Println(n)
	}
}
