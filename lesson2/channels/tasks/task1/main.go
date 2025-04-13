package main

import "fmt"

func main() {
	ch1 := make(chan int)
	go func() {
		ch1 <- 1
	}()

	fmt.Println(<-ch1)

	ch2 := make(chan int, 3)
	ch2 <- 1
	ch2 <- 2
	ch2 <- 3

	for i := 0; i < 3; i++ {
		fmt.Printf("len: %d, cap: %d val: %d\n", len(ch2), cap(ch2), <-ch2)
	}
}
