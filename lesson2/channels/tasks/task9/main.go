package main

import (
	"fmt"
	"sync"
	"time"
)

func Generate(wg *sync.WaitGroup, out chan<- int) {
	defer wg.Done()

	for i := 0; i < 100; i++ {
		out <- i
	}
	close(out)
}

func Square(wg *sync.WaitGroup, in <-chan int, out chan<- int) {
	defer wg.Done()

	for i := range in {
		out <- i * i
		time.Sleep(100 * time.Millisecond)
	}

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

	var wg2 sync.WaitGroup

	wg.Add(1)
	go Generate(&wg, ch1)

	for i := 0; i < 3; i++ {
		wg2.Add(1)
		go Square(&wg2, ch1, ch2)
	}

	go func() {
		wg2.Wait()
		close(ch2)
	}()

	go Print(&wg, ch2)
	wg.Wait()
}
