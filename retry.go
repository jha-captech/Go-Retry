package main

import (
	"context"
	"math"
	"math/rand"
	"time"
)

// Retry repeatedly calls the provided retryFunc until it succeeds or the maxDuration is exceeded.
// It uses an exponential backoff strategy for retries.
//
// Parameters:
// - ctx: The context to control the lifetime of the retry loop.
// - maxDuration: The maximum duration to keep retrying.
// - retryFunc: The function to call repeatedly until it succeeds.
//
// Returns:
// - error: The last error returned by retryFunc, or nil if retryFunc succeeds.
func Retry(ctx context.Context, maxDuration time.Duration, retryFunc func() error) error {
	_, err := RetryResult(ctx, maxDuration, func() (any, error) {
		return nil, retryFunc()
	})
	return err
}

// RetryResult repeatedly calls the provided retryFunc until it succeeds or the maxDuration is exceeded.
// It uses an exponential backoff strategy for retries.
//
// Parameters:
// - ctx: The context to control the lifetime of the retry loop.
// - maxDuration: The maximum duration to keep retrying.
// - retryFunc: The function to call repeatedly until it succeeds. It should return a result and an error.
//
// Returns:
// - T: The result returned by a successful call to retryFunc.
// - error: The last error returned by retryFunc, or nil if retryFunc succeeds.
func RetryResult[T any](ctx context.Context, maxDuration time.Duration, retryFunc func() (T, error)) (T, error) {
	var (
		returnData T
		err        error
	)
	const maxBackoffMilliseconds = 2_000.0

	ctx, cancelFunc := context.WithTimeout(ctx, maxDuration)
	defer cancelFunc()

	go func() {
		counter := 1.0
		for {
			counter++
			returnData, err = retryFunc()
			if err != nil {
				waitMilliseconds := math.Min(
					math.Pow(counter, 2)+float64(rand.Intn(10)),
					maxBackoffMilliseconds,
				)
				time.Sleep(time.Duration(waitMilliseconds) * time.Millisecond)
				continue
			}
			cancelFunc()
			return
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return returnData, err
		}
	}
}
