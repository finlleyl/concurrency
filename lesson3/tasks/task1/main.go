package main

import "fmt"

func main() {
	data := []int{1, 2, 3, 4, 5, 6, 7}

	inputCh := generator(data)

	consumer(inputCh)

}

func generator(data []int) <-chan int {
	inputCh := make(chan int)

	go func() {
		defer close(inputCh)

		for _, v := range data {
			inputCh <- v * v
		}
	}()

	return inputCh
}

func consumer(ch <-chan int) {
	for data := range ch {
		fmt.Println(data)
	}
}
