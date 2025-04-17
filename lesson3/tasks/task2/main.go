package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	ch := generator()

	// читаем данные из канала до его закрытия
	for v := range ch {
		fmt.Printf("Main received: %d\n", v)
	}

	fmt.Println("Main closed")
}

func generator() <-chan int {
	ch := make(chan int, 5)

	go func() {
		defer close(ch)
		for i := 0; i < 5; i++ {
			num := rand.Intn(1000)
			ch <- num
			fmt.Printf("Generator sent: %d\n", num)
			time.Sleep(time.Second)
		}
	}()

	fmt.Println("Generator started")
	return ch
}
