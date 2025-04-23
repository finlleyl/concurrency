package main

import "fmt"

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

func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	generatedCh := generator(nums...)
	result := square(generatedCh)

	for n := range result {
		fmt.Println(n)
	}
}
