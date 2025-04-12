package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	data := 0

	var wg sync.WaitGroup
	var mx sync.RWMutex

	// 5 читателей, удерживают RLock по 1 секунде
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("Reader %d trying to RLock at %v\n", id, time.Now())
			mx.RLock()
			fmt.Printf("Reader %d acquired RLock at %v\n", id, time.Now())
			fmt.Println(data)
			time.Sleep(1 * time.Second)
			mx.RUnlock()
			fmt.Printf("Reader %d released RLock at %v\n", id, time.Now())
		}(i)
	}

	// Дожидаемся, чтобы все RLock-и точно активировались
	time.Sleep(100 * time.Millisecond)

	// Запись — должна ждать завершения всех RLock-ов
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Printf("Writer trying to Lock at %v\n", time.Now())
		mx.Lock()
		fmt.Printf("Writer acquired Lock at %v\n", time.Now())
		data = 1
		time.Sleep(200 * time.Millisecond)
		mx.Unlock()
		fmt.Printf("Writer released Lock at %v\n", time.Now())
	}()

	wg.Wait()
}
