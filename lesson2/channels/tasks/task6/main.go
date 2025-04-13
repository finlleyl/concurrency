package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string)
	go func() {
		time.Sleep(3 * time.Second)
		ch <- "Hello"
	}()

	select {
	case <-ch:
		fmt.Print("World")

	case <-time.After(2 * time.Second):
		fmt.Println("Timeout")
		return

	}
}
