package main

import (
	"fmt"
	"log"
	"math/rand"
)

func main() {
	var data [100]int
	for i := 0; i < 100; i++ {
		data[i] = rand.Intn(100)
	}

	handler(data[:])
}

func handler(data []int) {
	doneCh := make(chan struct{})
	defer close(doneCh)

	inputCh := generator(data, doneCh)

	for v := range inputCh {
		if v == 42 {
			fmt.Println("❌ Unexpected int: 42 — останавливаем обработку")
			return
		}
		fmt.Println("✅", v)
	}

	fmt.Println("🎉 Обработка завершена без остановки")
}

func generator(data []int, doneCh <-chan struct{}) <-chan int {
	inputCh := make(chan int)

	go func() {
		defer func() {
			log.Println("📦 Канал inputCh закрыт (генератор завершён)")
			close(inputCh)
		}()

		for _, v := range data {
			select {
			case inputCh <- v:
				log.Printf("🔄 Отправлено: %d\n", v)
			case <-doneCh:
				log.Println("🛑 Generator остановлен по сигналу doneCh")
				return
			}
		}
	}()

	return inputCh
}
