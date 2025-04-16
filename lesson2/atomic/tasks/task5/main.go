package main

import (
	"fmt"
	"sync/atomic"
)

func main() {
	var x int64 = 10

	res := atomic.CompareAndSwapInt64(&x, 5, 64)

	fmt.Printf("Results of CompareAndSwapInt64: %v\nNew value: %d", res, atomic.LoadInt64(&x))
}
