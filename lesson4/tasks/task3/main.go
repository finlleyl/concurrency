package main

import "fmt"

type Stage func(<-chan struct{}, <-chan int) <-chan int

func generator(doneCh <-chan struct{}, nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			select {
			case out <- n:
			case <-doneCh:
				return
			}
		}
	}()
	return out
}

func square(doneCh <-chan struct{}, in <-chan int) <-chan int {
	res := make(chan int)
	go func() {
		defer close(res)
		for n := range in {
			select {
			case res <- n * n:
			case <-doneCh:
				return
			}
		}
	}()
	return res
}

func increment(doneCh <-chan struct{}, in <-chan int) <-chan int {
	res := make(chan int)
	go func() {
		defer close(res)
		for n := range in {
			select {
			case res <- n + 1:
			case <-doneCh:
				return
			}
		}
	}()
	return res
}

func pipeline(doneCh <-chan struct{}, in <-chan int, stages ...Stage) <-chan int {
	ch := in
	for _, stage := range stages {
		ch = stage(doneCh, ch)
	}
	return ch
}

func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	doneCh := make(chan struct{})

	out := pipeline(doneCh,
		generator(doneCh, nums...),
		increment,
		square,
	)

	for n := range out {
		if n == 5 {
			close(doneCh) // сигналим об отмене всем этапам
			break         // прекращаем чтение из out
		}
		fmt.Println(n)
	}
}
