package hw05parallelexecution

import (
	"errors"
	"sync"
)

type Task func() error

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	mu                     sync.RWMutex
)

func Run(tasks []Task, n, m int) error {
	wg := sync.WaitGroup{}

	taskQueue := make(chan Task, n)
	currentErrors := 0
	terminateSignal := make(chan struct{}, 1)

	// Producer
	wg.Add(1)
	go func() {
		wg.Done()
		producer(taskQueue, tasks, terminateSignal)
	}()

	// Consumers
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			consumer(taskQueue, &currentErrors, &m, terminateSignal)
			wg.Done()
		}()
	}

	wg.Wait()

	if currentErrors >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func producer(taskQueue chan<- Task, tasks []Task, terminateSignal <-chan struct{}) {
L:
	for {
		select {
		case taskQueue <- tasks[0]:
			if len(tasks) == 1 {
				break L
			}
			tasks = tasks[1:]
		case <-terminateSignal:
			break L
		}
	}

	close(taskQueue)
}

func consumer(taskQueue <-chan Task, currentErrors *int, maxErrors *int, terminateSignal chan<- struct{}) {
	for task := range taskQueue {
		if task() != nil {
			mu.Lock()
			*(currentErrors)++
			errorState := *(currentErrors) >= *(maxErrors)
			mu.Unlock()

			if errorState {
				select {
				case terminateSignal <- struct{}{}:
				default:
				}
				break
			}
		}
	}
}
