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
	// Place your code here.
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	errorCounter := int64(0)
	wg := sync.WaitGroup{}
	channel := make(chan Task)

	for i := 0; i < n; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for task := range channel {
				err := task()
				if err != nil {
					atomic.AddInt64(&errorCounter, 1)
				}
			}
		}()
	}

	for _, task := range tasks {
		if atomic.LoadInt64(&errorCounter) >= int64(m) {
			break
		}
		channel <- task
	}
	close(channel)

	wg.Wait()

	if atomic.LoadInt64(&errorCounter) >= int64(m) {
		return ErrErrorsLimitExceeded
	}
	return nil
}
