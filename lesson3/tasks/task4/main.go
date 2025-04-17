package main

import (
	"fmt"
	"golang.org/x/sync/errgroup"
)

func main() {
	g := new(errgroup.Group)

	inputCh := generator([]int{1, 2, 3, 4, 5, 6, 7})

	for data := range inputCh {
		data := data
		g.Go(func() error {
			if data == 5 {
				return fmt.Errorf("data is 5")
			}
			fmt.Printf("Data: %d\n", data)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Println(err)
	}

}

func generator(data []int) <-chan int {
	ch := make(chan int)

	go func() {
		defer close(ch)
		for _, v := range data {
			ch <- v
		}
	}()

	return ch
}
