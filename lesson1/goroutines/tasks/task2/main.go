package main

import (
	"fmt"
	"time"
)

func printHello() {
	for i := 0; i < 5; i++ {
		fmt.Println("Hello")
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	go printHello()

	go func() {
		for i := 0; i < 5; i++ {
			fmt.Println("Main still working...")
			time.Sleep(150 * time.Millisecond)
		}
	}()

	time.Sleep(1 * time.Second)
}
