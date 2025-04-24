package main

import (
	"fmt"
	"time"
)

// worker берёт задачи из jobs и пишет результат в results
func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Printf("Worker %d: got job %d\n", id, j)
		time.Sleep(200 * time.Millisecond) // имитируем работу
		results <- j * 2
	}
}

func main() {
	jobs := make(chan int, 10)
	results := make(chan int, 10)

	// запускаем 4 воркера
	for w := 1; w <= 4; w++ {
		go worker(w, jobs, results)
	}

	// отправляем задачи
	for j := 1; j <= 10; j++ {
		jobs <- j
	}
	close(jobs) // сигнализируем воркерам, что задач больше не будет

	// читаем результаты
	for i := 0; i < 10; i++ {
		fmt.Println("Result:", <-results)
	}
}
