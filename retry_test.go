package main

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type retryFuncMock struct {
	callCount     int
	shouldSucceed bool
}

func (r *retryFuncMock) call() error {
	r.callCount++
	if r.shouldSucceed && r.callCount > 3 {
		return nil
	}
	return errors.New("mock error")
}

func (r *retryFuncMock) callWithResult() (string, error) {
	r.callCount++
	if r.shouldSucceed && r.callCount > 3 {
		return "success", nil
	}
	return "", errors.New("mock error")
}

func TestRetry(t *testing.T) {
	testCases := map[string]struct {
		maxDuration   time.Duration
		shouldSucceed bool
		expectError   bool
	}{
		"Success after retries": {
			maxDuration:   5 * time.Second,
			shouldSucceed: true,
			expectError:   false,
		},
		"Failure after max duration": {
			maxDuration:   1 * time.Second,
			shouldSucceed: false,
			expectError:   true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mockFunc := &retryFuncMock{shouldSucceed: tc.shouldSucceed}

			err := Retry(context.Background(), tc.maxDuration, mockFunc.call)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRetryResult(t *testing.T) {
	testCases := map[string]struct {
		maxDuration   time.Duration
		shouldSucceed bool
		expectError   bool
		expectResult  string
	}{
		"Success after retries": {
			maxDuration:   5 * time.Second,
			shouldSucceed: true,
			expectError:   false,
			expectResult:  "success",
		},
		"Failure after max duration": {
			maxDuration:   1 * time.Second,
			shouldSucceed: false,
			expectError:   true,
			expectResult:  "",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mockFunc := &retryFuncMock{shouldSucceed: tc.shouldSucceed}

			result, err := RetryResult(context.Background(), tc.maxDuration, mockFunc.callWithResult)

			if tc.expectError {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectResult, result)
			}
		})
	}
}
