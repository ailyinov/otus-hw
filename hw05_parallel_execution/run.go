package hw05parallelexecution

import (
	"context"
	"errors"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func worker(jobs <-chan Task, results chan<- error, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			for j := range jobs {
				results <- j()
			}
		}
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	jobs := make(chan Task, n)
	results := make(chan error, len(tasks))
	defer close(jobs)

	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	for w := 0; w < n; w++ {
		go worker(jobs, results, ctx)
	}

	for _, t := range tasks {
		jobs <- t
	}

	errsCount := 0
	for range tasks {
		if err := <-results; err != nil {
			if m <= 0 {
				continue
			}
			if errsCount++; errsCount >= m {
				return ErrErrorsLimitExceeded
			}
		}
	}
	return nil
}
