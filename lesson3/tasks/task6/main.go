package main

import (
	"context"
	"log"
	"math/rand"
)

func main() {
	data := make([]int, 100)
	for i := 0; i < 100; i++ {
		data[i] = rand.Intn(105)
	}

	handler(data)
}

func handler(data []int) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	inputCh := generator(data, ctx)

	for v := range inputCh {
		if v > 100 {
			log.Printf("❌ Unexpected int: %d — вызываем cancel()\n", v)
			cancel()
			return
		}
		log.Printf("✅ Принято: %d\n", v)
	}
	log.Println("🎉 Обработка завершена без отмены")
}

func generator(data []int, ctx context.Context) <-chan int {
	inputCh := make(chan int)

	go func() {
		defer close(inputCh)

		for _, v := range data {
			select {
			case <-ctx.Done():
				log.Println("Generator Done (early)")
				return
			default:
			}

			select {
			case inputCh <- v:
			case <-ctx.Done():
				log.Println("generator done (during send)")
			}
		}
	}()

	return inputCh
}
