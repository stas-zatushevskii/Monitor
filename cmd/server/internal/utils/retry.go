package utils

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/models"
)

type DBError struct {
	Msg  string
	Code int
}

func (e DBError) Error() string {
	return e.Msg
}

func (e DBError) Retryable() bool {
	return e.Code == 1205 || e.Code == 57014 // SQL timeout, cancel
}

type RetryableError interface {
	error
	Retryable() bool
}

// RetryWithContext repeats the request with context until the limit of access is exceeded
// default limit is 5, default delay is 1 second
// time gap between attempts scales according formula
//
//	delay += 2 * time.Second
//
// function will stop if error from request is not retryable
func RetryWithContext(
	ctx context.Context,
	fn func(ctx context.Context, data []models.Metrics) error,
	data []models.Metrics,
) error {
	var err error
	delay := 1 * time.Second
	retries := 5

	err = fn(ctx, data)
	if err == nil || !isRetryable(err) {
		return err
	}

	for i := 0; i <= retries; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
			fmt.Printf("Retryable error: %v. Retrying attempt %d/%d...\n", err, i+1, retries)
			err = fn(ctx, data)
			if err == nil || !isRetryable(err) {
				return err
			}
			delay += 2 * time.Second
		}
	}

	return err
}

// RetryGetDataByName repeats the request until the limit of access is exceeded
// default limit is 5, default delay is 1 second
// time gap between attempts scales according formula
//
//	delay += 2 * time.Second
//
// function will stop if error from request is not retryable
func RetryGetDataByName(
	fn func(nameMetric, typeMetric string) (string, error),
	nameMetric, typeMetric string,
) (string, error) {
	retries := 5
	result, err := fn(nameMetric, typeMetric)
	if err == nil || !isRetryable(err) {
		return result, err
	}

	delay := 1 * time.Second

	for i := 0; i <= retries; i++ {
		fmt.Printf("Retryable error: %v. Retrying attempt %d/%d...\n", err, i+1, retries)
		time.Sleep(delay)
		result, err = fn(nameMetric, typeMetric)

		if err == nil || !isRetryable(err) {
			return result, err
		}

		delay += 2 * time.Second
	}

	return result, err
}

// RetrySetJSONData repeats the request until the limit of access is exceeded
// default limit is 5, default delay is 1 second
// time gap between attempts scales according formula
//
//	delay += 2 * time.Second
//
// function will stop if error from request is not retryable
func RetrySetJSONData(
	fn func(data models.Metrics) error,
	data models.Metrics,
) error {
	err := fn(data)
	if err == nil || !isRetryable(err) {
		return err
	}
	retries := 5
	delay := 1 * time.Second

	for i := 0; i <= retries; i++ {
		fmt.Printf("Retryable error: %v. Retrying attempt %d/%d...\n", err, i+1, retries)
		time.Sleep(delay)
		err = fn(data)

		if err == nil || !isRetryable(err) {
			return err
		}
		delay += 2 * time.Second
	}

	return err
}

// RetrySetURLData repeats the request until the limit of access is exceeded
// default limit is 5, default delay is 1 second
// time gap between attempts scales according formula
//
//	delay += 2 * time.Second
//
// function will stop if error from request is not retryable
func RetrySetURLData(
	fn func(nameMetric, dataMetric, typeMetric string) error,
	nameMetric, dataMetric, typeMetric string,
) error {
	err := fn(nameMetric, dataMetric, typeMetric)
	if err == nil || !isRetryable(err) {
		return err
	}
	retries := 5
	delay := 1 * time.Second

	for i := 0; i <= retries; i++ {
		fmt.Printf("Retryable error: %v. Retrying attempt %d/%d...\n", err, i+1, retries)
		time.Sleep(delay)
		err = fn(nameMetric, dataMetric, typeMetric)

		if err == nil || !isRetryable(err) {
			return err
		}
		delay += 2 * time.Second
	}

	return err
}

// RetryGetAllGaugeMetrics repeats the request until the limit of access is exceeded
// default limit is 5, default delay is 1 second
// time gap between attempts scales according formula
//
//	delay += 2 * time.Second
//
// function will stop if error from request is not retryable
func RetryGetAllGaugeMetrics(
	fn func() (map[string]float64, error),

) (map[string]float64, error) {
	result, err := fn()
	if err == nil || !isRetryable(err) {
		return result, err
	}

	retries := 5
	delay := 1 * time.Second

	for i := 0; i <= retries; i++ {
		fmt.Printf("Retryable error: %v. Retrying gauge metrics attempt %d/%d...\n", err, i+1, retries)
		time.Sleep(delay)
		result, err = fn()
		if err == nil || !isRetryable(err) {
			return result, err
		}
		delay += 2 * time.Second
	}

	return result, err
}

// RetryGetAllCounterMetrics repeats the request until the limit of access is exceeded
// default limit is 5, default delay is 1 second
// time gap between attempts scales according formula
//
//	delay += 2 * time.Second
//
// function will stop if error from request is not retryable
func RetryGetAllCounterMetrics(
	fn func() (map[string]int64, error),
) (map[string]int64, error) {
	result, err := fn()
	if err == nil || !isRetryable(err) {
		return result, err
	}

	retries := 5
	delay := 1 * time.Second

	for i := 0; i <= retries; i++ {
		fmt.Printf("Retryable error: %v. Retrying counter metrics attempt %d/%d...\n", err, i+1, retries)
		time.Sleep(delay)
		result, err = fn()
		if err == nil || !isRetryable(err) {
			return result, err
		}
		delay += 2 * time.Second
	}

	return result, err
}

// isRetryable checks if error got from request is retryable
//
//	request is Timeout from bd and Cancel
func isRetryable(err error) bool {
	if err == nil {
		return false
	}
	var r RetryableError
	if errors.As(err, &r) {
		return r.Retryable()
	}
	return false
}
