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
	logWithNumber := make([]string, len(e))
	for i, l := range e {
		if l != nil {
			logWithNumber[i] = fmt.Sprintf("#%d: %s", i+1, l.Error())
		}
	}

	return fmt.Sprintf("Retry attempts failed:\n%s", strings.Join(logWithNumber, "\n"))
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
		time.Sleep(*attemptTimeout)

		err = retryableFunc()
		if err != nil {
			errorSlice = append(errorSlice, err)

			if shouldRetryFunc(err) {
				attemptTimeout = nextAttemptTimeout(attemptTimeout)
			} else {
				return errorSlice
			}
		}

		curAttempt++
	}

	return errorSlice
}

func nextAttemptTimeout(prevT *time.Duration) *time.Duration {
	if prevT == nil {
		t := time.Second
		return &t
	}
	newT := *prevT + 2*time.Second
	return &newT
}
