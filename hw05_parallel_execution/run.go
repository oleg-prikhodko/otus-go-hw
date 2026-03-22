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

	quitChan := make(chan struct{})
	closer := func() { close(quitChan) }
	var once sync.Once
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
				select {
				case <-quitChan:
					return
				default:
				}


				if curTaskIdx >= len(tasks) {
					return
				}

				taskMutex.Lock()
				task := tasks[curTaskIdx]
				curTaskIdx++
				taskMutex.Unlock()

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

	wg.Wait()

	if errCount >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}
