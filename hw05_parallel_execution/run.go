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

	tasksChan := make(chan Task, len(tasks))
	quitChan := make(chan struct{})
	closer := func() { close(quitChan) }
	var once sync.Once
	var errCount int
	var errMutex sync.Mutex
	var wg sync.WaitGroup

	wg.Add(n)
	for i := 0; i < n; i++ {
		// worker goroutine
		go func() {
			defer wg.Done()
			for {
				select {
				case <-quitChan:
					return
				default:
				}

				task, moreTasks := <-tasksChan
				if !moreTasks {
					return
				}
				if err := task(); err != nil {
					errMutex.Lock()
					if errCount++; errCount >= m {
						once.Do(closer)
					}
					errMutex.Unlock()
				}
			}
		}()
	}

	for i := 0; i < len(tasks); i++ {
		tasksChan <- tasks[i]
	}
	close(tasksChan)

	wg.Wait()

	if errCount >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}
