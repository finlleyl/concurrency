package main

import (
	"fmt"
	"sync/atomic"
)

func main() {
	var x int64 = 10

	fmt.Println(atomic.SwapInt64(&x, 99))
}
