package main

import (
	"errors"
	"time"
)

// retry will retry a given function n times with a wait of a given duration between each retry attempt.
func retry(retryCount int, waitTime time.Duration, fn func() error) error {
	_, err := retryWithReturn(retryCount, waitTime, func() (any, error) {
		return nil, fn()
	})
	return err
}

// retryWithReturn will retry a given function n times with a wait of a given duration between each
// retry attempt. retryWithReturn is intended for functions where a return values is needed.
func retryWithReturn[T any](retryCount int, waitTime time.Duration, fn func() (T, error)) (T, error) {
	if retryCount < 1 {
		return *new(T), errors.New("retryCount of less than 1 is not permitted")
	}
	for i := 0; i < retryCount; i++ {
		t, err := fn()
		if err != nil {
			if i == retryCount-1 {
				return t, err
			}
			time.Sleep(waitTime)
			continue
		}
		return t, nil
	}
	return *new(T), errors.New("default return reached in retryWithReturn")
}
