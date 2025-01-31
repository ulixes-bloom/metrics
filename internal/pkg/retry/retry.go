// Package retry provides a mechanism to perform retries for operations that can fail,
// allowing for custom retry logic and handling multiple retry attempts with backoff.
package retry

import (
	"fmt"
	"strings"
	"time"
)

// RetryableFunc defines a function type that represents an operation that can be retried.
// It returns an error if the operation fails and nil if it succeeds.
type RetryableFunc func() error

// ShouldRetryFunc defines a function type that determines whether a failed operation
// should be retried based on the error returned by the RetryableFunc.
type ShouldRetryFunc func(error) bool

// ErrorSlice is a slice of errors used to store the errors encountered during retry attempts.
// It implements the Error interface for custom error handling.
type ErrorSlice []error

// Error formats the ErrorSlice to display all the errors encountered during retry attempts
// in a human-readable format, including the index of each failed attempt.
func (e ErrorSlice) Error() string {
	logWithNumber := []string{}
	for i, l := range e {
		if l != nil {
			logWithNumber = append(logWithNumber, fmt.Sprintf("#%d: %s", i+1, l.Error()))
		}
	}

	return fmt.Sprintf("Retry attempts failed: %s", strings.Join(logWithNumber, "|"))
}

// Do performs the given retryable function with the specified retry logic.
// It retries the operation up to the specified number of attempts if the error returned
// from the retryable function satisfies the shouldRetry function.
func Do(retryableFunc RetryableFunc, shouldRetryFunc ShouldRetryFunc, attempts uint) error {
	err := retryableFunc()
	if err == nil {
		return nil
	}
	shouldRetry := shouldRetryFunc(err)
	if !shouldRetry {
		return fmt.Errorf("retry.do: %w", err)
	}

	errorSlice := ErrorSlice{err}
	attemptTimeout := nextAttemptTimeout(nil)
	curAttempt := uint(2)

	for shouldRetry && curAttempt != attempts {
		time.Sleep(attemptTimeout)

		err = retryableFunc()
		if err != nil {
			errorSlice = append(errorSlice, err)

			if shouldRetryFunc(err) {
				attemptTimeout = nextAttemptTimeout(&attemptTimeout)
			} else {
				shouldRetry = false
			}
		}

		curAttempt++
	}

	return errorSlice
}

// nextAttemptTimeout calculates the timeout for the next retry attempt based on the previous timeout.
// If it's the first retry attempt, it returns a default timeout of 1 second. Otherwise, it increases
// the timeout by 2 seconds for each subsequent retry.
func nextAttemptTimeout(prevT *time.Duration) time.Duration {
	if prevT == nil {
		t := time.Second
		return t
	}
	newT := *prevT + 2*time.Second
	return newT
}
