package main

import (
	"fmt"
)

type Result struct {
	Data int
	Err  error
}

func main() {
	inputCh := generator([]int{1, 2, 3, 4, 5, 6, 7})

	outputCh := consumer(inputCh)

	for res := range outputCh {
		if res.Err != nil {
			fmt.Println(res.Err)
		} else {
			fmt.Println(res.Data)
		}
	}
}

func generator(data []int) <-chan int {
	inputCh := make(chan int)

	go func() {
		defer close(inputCh)

		for _, v := range data {
			inputCh <- v
		}
	}()

	return inputCh
}

func consumer(ch <-chan int) <-chan Result {
	outputCh := make(chan Result)

	go func() {
		defer close(outputCh)
		for data := range ch {
			res := Result{data, nil}
			if data%2 == 0 {
				res.Err = fmt.Errorf("even number: %d", data)
			}

			outputCh <- res
		}
	}()
	return outputCh
}
