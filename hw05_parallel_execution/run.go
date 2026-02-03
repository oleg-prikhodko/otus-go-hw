package hw05parallelexecution

import (
	"errors"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	t := make(chan Task, len(tasks))
	e := make(chan error, m)

	for i := 0; i < n; i++ {
		go worker(t, e)
	}
	for i := 0; i < len(tasks); i++ {
		t <- tasks[i]
	}
	close(t)

	for i := 0; ; {
		if i >= m {
			// todo: stop goroutines
			return ErrErrorsLimitExceeded
		}
		err := <-e
		if err != nil {
			i++
		}
	}

	// Place your code here.
	return nil
}

func worker(tasks <-chan Task, errors chan<- error) {
	for t := range tasks {
		errors <- t()
	}
}
