package main

import (
	"sync"
)

func main() {
	var mu sync.Mutex

	mu.Unlock() // ← паника!
}

// fatal error: sync: unlock of unlocked mutex
