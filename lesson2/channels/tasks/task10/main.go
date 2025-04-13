package main

import (
	"fmt"
	"runtime"
	"time"
)

type Task struct {
	id int
}

type Queue struct {
	ch chan *Task
}

func NewQueue() *Queue {
	ch := make(chan *Task, 1)
	return &Queue{ch}
}

func (q *Queue) Push(t *Task) {
	q.ch <- t
}

func (q *Queue) Pop() *Task {
	return <-q.ch
}

type Worker struct {
	id    int
	queue *Queue
}

func NewWorker(id int, queue *Queue) *Worker {
	return &Worker{id: id, queue: queue}
}

func (w *Worker) Run() {
	for {
		task := w.queue.Pop()
		fmt.Printf("Worker %d got task %d\n", w.id, task.id)
		time.Sleep(time.Second)
		fmt.Printf("Worker %d finished task %d-->%d\n", w.id, task.id, task.id*task.id)
	}
}

func main() {
	queue := NewQueue()

	for i := 0; i < runtime.NumCPU(); i++ {
		w := NewWorker(i, queue)
		go w.Run()
	}

	for i := 0; i < 500; i++ {
		task := &Task{id: i}
		queue.Push(task)
	}
}
