package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

var ErrWorkerCountIncorrect = errors.New("n must be > 0")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if n <= 0 {
		return ErrWorkerCountIncorrect
	}
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	var errCount int
	var errMutex sync.Mutex
	var wg sync.WaitGroup

	var taskMutex sync.Mutex
	var curTaskIdx int

	wg.Add(n)
	for i := 0; i < n; i++ {
		// worker goroutine
		go func() {
			defer wg.Done()
			for {
				taskMutex.Lock()
				if curTaskIdx >= len(tasks) {
					taskMutex.Unlock()
					return
				}
				task := tasks[curTaskIdx]
				curTaskIdx++
				taskMutex.Unlock()

				if err := task(); err != nil {
					errMutex.Lock()
					errCount++
					exceeded := errCount >= m
					errMutex.Unlock()
					if exceeded {
						return
					}
				}
			}
		}()
	}

	wg.Wait()

	if errCount >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}
