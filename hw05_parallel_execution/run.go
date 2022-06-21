package hw05parallelexecution

import (
	"errors"
	"sync"
)

type Task func() error

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
var ErrZeroConsumersCount = errors.New("zero consumers count")
var mu sync.RWMutex
var currentErrorsCount, maxErrorsCount int
var terminateSignal chan struct{}
var errorState bool

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	// Init
	errorState = false
	if m == 0 {
		return ErrErrorsLimitExceeded
	}
	if n == 0 {
		return ErrZeroConsumersCount
	}

	maxErrorsCount = m
	terminateSignal = make(chan struct{}, 1)
	defer close(terminateSignal)

	wg := sync.WaitGroup{}

	// Init producer
	queue := make(chan Task, n)

	wg.Add(1)
	go func() {
		defer wg.Done()
		producer(queue, tasks)
	}()

	// Init consumers
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			consumer(queue)
		}()
	}

	// Wait until done
	wg.Wait()

	if errorState {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func producer(queue chan<- Task, tasks []Task) {
L:
	for {
		// Priority
		select {
		case <-terminateSignal:
			break L

		default:
		}

		// Trying push task to queue
		task := tasks[0]
		if task == nil {
			break L
		}

		select {
		case queue <- task:
			if len(tasks) > 1 {
				tasks = tasks[1:]
			} else {
				break L
			}
		case <-terminateSignal:
			break L
		default:
		}
	}

	close(queue)
}

func consumer(queue <-chan Task) {
	for task := range queue {
		err := task()

		if err != nil && checkErrorState() {
			return
		}
	}
}

func checkErrorState() bool {
	mu.Lock()
	currentErrorsCount++
	mu.Unlock()

	mu.RLock()
	tooManyErrors := currentErrorsCount >= maxErrorsCount
	mu.RUnlock()

	if tooManyErrors {
		mu.Lock()
		errorState = true
		mu.Unlock()

		// Non blocking
		select {
		case terminateSignal <- struct{}{}:
		default:
		}

		return true
	}

	return false
}
