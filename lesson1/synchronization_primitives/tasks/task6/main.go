package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	mx := sync.Mutex{}
	cond := sync.NewCond(&mx)

	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("Сигнал!")
		cond.Signal()
	}()

	mx.Lock()
	fmt.Println("Жду сигнала")
	cond.Wait()
	fmt.Println("Получил сигнал")
	mx.Unlock()
}
