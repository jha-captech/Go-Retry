package main

import "time"

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
	var (
		err error
		t   T
	)
	for i := 0; i < retryCount; i++ {
		t, err = fn()
		if err != nil {
			if i >= retryCount {
				return t, err
			}
			time.Sleep(waitTime)
			continue
		}
		return t, nil
	}
	return t, err
}
