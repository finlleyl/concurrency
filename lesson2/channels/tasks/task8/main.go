package main

import (
	"fmt"
	"sync"
	"time"
)

func Generate(wg *sync.WaitGroup, out chan<- int) {
	defer wg.Done()

	for i := 0; i < 10; i++ {
		out <- i
		time.Sleep(100 * time.Millisecond)
	}
	close(out)
}

func Square(wg *sync.WaitGroup, in <-chan int, out chan<- int) {
	defer wg.Done()

	for i := range in {
		out <- i * i
		time.Sleep(100 * time.Millisecond)
	}
	close(out)
}

func Print(wg *sync.WaitGroup, in <-chan int) {
	defer wg.Done()

	for i := range in {
		fmt.Println(i)
	}

}
func main() {

	ch1 := make(chan int)
	ch2 := make(chan int)
	var wg sync.WaitGroup

	wg.Add(3)
	go Generate(&wg, ch1)
	go Square(&wg, ch1, ch2)
	go Print(&wg, ch2)
	wg.Wait()
}
