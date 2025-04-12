package main

import (
	"fmt"
	"time"
)

func main() {
	for i := 0; i < 10; i++ {
		go func(n int) {
			fmt.Println(n)
		}(i)
	}

	time.Sleep(100 * time.Millisecond)
}
