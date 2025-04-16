package main

import (
	"fmt"
	"sync/atomic"
)

func main() {
	var x int64 = 99

	for {
		old := atomic.LoadInt64(&x)

		if old >= 100 {
			fmt.Println("x >= 100")
			break
		}

		if atomic.CompareAndSwapInt64(&x, old, 100) {
			fmt.Println("x was ", old, " and is now 100")
			break
		}
	}
}
