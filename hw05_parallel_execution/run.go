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
	wg := sync.WaitGroup{}
	taskChannel := make(chan Task, len(tasks))
	var errCount int64

	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()

			for task := range taskChannel {
				err := task()

				if atomic.LoadInt64(&errCount) >= int64(m) {
					return
				}

				if err != nil {
					atomic.AddInt64(&errCount, 1)
				}
			}
		}()
	}

	for _, task := range tasks {
		taskChannel <- task
	}

	close(taskChannel)

	wg.Wait()

	if atomic.LoadInt64(&errCount) >= int64(m) {
		return ErrErrorsLimitExceeded
	}

	return nil
}
