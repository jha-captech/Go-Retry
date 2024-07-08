package main

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRetry(t *testing.T) {
	tests := map[string]struct {
		retryCount int
		waitTime   time.Duration
		fn         func() error
		expectErr  bool
	}{
		"Success on first attempt": {
			retryCount: 3,
			waitTime:   100 * time.Millisecond,
			fn: func() error {
				return nil
			},
			expectErr: false,
		},
		"Success on second attempt": {
			retryCount: 3,
			waitTime:   100 * time.Millisecond,
			fn: func() func() error {
				attempts := 0
				return func() error {
					attempts++
					if attempts < 2 {
						return errors.New("error")
					}
					return nil
				}
			}(),
			expectErr: false,
		},
		"All attempts fail": {
			retryCount: 3,
			waitTime:   100 * time.Millisecond,
			fn: func() error {
				return errors.New("error")
			},
			expectErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := retry(tt.retryCount, tt.waitTime, tt.fn)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRetryWithReturn(t *testing.T) {
	tests := map[string]struct {
		retryCount int
		waitTime   time.Duration
		fn         func() (int, error)
		expectVal  int
		expectErr  bool
	}{
		"Success on first attempt": {
			retryCount: 3,
			waitTime:   100 * time.Millisecond,
			fn: func() (int, error) {
				return 42, nil
			},
			expectVal: 42,
			expectErr: false,
		},
		"Success on second attempt": {
			retryCount: 3,
			waitTime:   100 * time.Millisecond,
			fn: func() func() (int, error) {
				attempts := 0
				return func() (int, error) {
					attempts++
					if attempts < 2 {
						return 0, errors.New("error")
					}
					return 42, nil
				}
			}(),
			expectVal: 42,
			expectErr: false,
		},
		"All attempts fail": {
			retryCount: 3,
			waitTime:   100 * time.Millisecond,
			fn: func() (int, error) {
				return 0, errors.New("error")
			},
			expectVal: 0,
			expectErr: true,
		},
		"Check last return error not reached": {
			retryCount: 0,
			waitTime:   100 * time.Millisecond,
			fn: func() (int, error) {
				return 0, errors.New("error")
			},
			expectVal: 0,
			expectErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			val, err := retryWithReturn(tt.retryCount, tt.waitTime, tt.fn)
			if tt.expectErr {
				assert.Error(t, err)
				if tt.retryCount == 0 {
					assert.Equal(t, "default return reached in retryWithReturn", err.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectVal, val)
			}
		})
	}
}
