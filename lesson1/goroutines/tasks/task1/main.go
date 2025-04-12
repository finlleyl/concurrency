package main

import (
	"fmt"
	"time"
)

func main() {
	go func() {
		fmt.Println("Горутина 1")
	}()
	go func() {
		fmt.Println("Горутина 2")
	}()

	time.Sleep(100 * time.Millisecond)
}
