package utils

import (
	"fmt"
	"time"

	"github.com/stas-zatushevskii/Monitor/cmd/agent/internal/types"
)

func RetryRequest[T types.MetricData](
	fn func(m T, url, hashKey string) error,
	data T,
	url, hashKey string,
) error {
	retryCount := 3
	timeout := 1

	err := fn(data, url, hashKey)
	if err == nil || !isRetryable(err) {
		return err
	}

	for i := 0; i < retryCount; i++ {
		fmt.Printf("Retryable error: %v. Retrying in %d seconds...\n", err, timeout)
		time.Sleep(time.Duration(timeout) * time.Second)
		timeout += 2
		err = fn(data, url, hashKey)

		if err == nil || !isRetryable(err) {
			break
		}
	}

	return err
}

func isRetryable(err error) bool {
	if err == nil {
		return false
	}
	if r, ok := err.(interface{ Retryable() bool }); ok {
		return r.Retryable()
	}
	return false
}
