package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"runtime"
	"time"
)

type Job struct {
	Num int
}

type Result struct {
	Num int
}

func worker(id int, ctx context.Context, jobs <-chan Job, results chan<- Result) error {
	for job := range jobs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if job.Num == 90_000 {
			return errors.New("Unresolved int")
		}
		fmt.Printf("Worker %d received job %d\n", id, job)
		time.Sleep(time.Second)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case results <- Result{Num: job.Num * job.Num}:
			fmt.Printf("Worker %d completed job %d\n", id, job)
		}
	}

	return nil
}

func main() {
	g, ctx := errgroup.WithContext(context.Background())
	jobs := make(chan Job)
	results := make(chan Result)

	go func() {
		if err := g.Wait(); err != nil {
			fmt.Println(err)
		}
		close(jobs)
	}()

	go func() {
		for res := range results {
			fmt.Printf("Main recived res %d\n", res.Num)
		}
	}()

	for w := 0; w < runtime.NumCPU(); w++ {
		w := w
		g.Go(func() error {
			return worker(w, ctx, jobs, results)
		})
	}

	go func() {
		for j := 0; j < 10_000; j++ {
			select {
			case <-ctx.Done():
				break
			case jobs <- Job{Num: j}:
			}
		}
		close(jobs)
	}()

	<-ctx.Done()
}
