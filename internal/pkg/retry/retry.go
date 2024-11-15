package retry

import (
	"fmt"
	"strings"
	"time"
)

type RetryableFunc func() error

type ShouldRetryFunc func(error) bool

type ErrorSlice []error

func (e ErrorSlice) Error() string {
	logWithNumber := []string{}
	for i, l := range e {
		if l != nil {
			logWithNumber = append(logWithNumber, fmt.Sprintf("#%d: %s", i+1, l.Error()))
		}
	}

	return fmt.Sprintf("Retry attempts failed: %s", strings.Join(logWithNumber, "|"))
}

func Do(retryableFunc RetryableFunc, shouldRetryFunc ShouldRetryFunc, attempts uint) error {
	err := retryableFunc()
	if err == nil {
		return nil
	}
	shouldRetry := shouldRetryFunc(err)
	if !shouldRetry {
		return err
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

func nextAttemptTimeout(prevT *time.Duration) time.Duration {
	if prevT == nil {
		t := time.Second
		return t
	}
	newT := *prevT + 2*time.Second
	return newT
}
