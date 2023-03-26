package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var cntErr int32
	taskCh := make(chan Task)
	wg := sync.WaitGroup{}
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for task := range taskCh {
				if err := task(); err != nil {
					atomic.AddInt32(&cntErr, 1)
				}
			}
		}()
	}

	for _, task := range tasks {
		if atomic.LoadInt32(&cntErr) >= int32(m) {
			break
		}
		taskCh <- task
	}

	close(taskCh)
	wg.Wait()

	if cntErr >= int32(m) {
		return ErrErrorsLimitExceeded
	}

	return nil
}
