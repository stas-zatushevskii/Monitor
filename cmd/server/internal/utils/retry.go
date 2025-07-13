package utils

import (
	"context"
	"fmt"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/models"
	"time"
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

func RetryWithContext(
	ctx context.Context,
	fn func(ctx context.Context, data []models.Metrics) error,
	data []models.Metrics,
) error {
	var err error
	delay := 1 * time.Second
	retries := 100

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

func RetryGetDataByName(
	fn func(nameMetric, typeMetric string) (string, error),
	nameMetric, typeMetric string,
) (string, error) {
	retries := 100
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

func RetrySetJSONData(
	fn func(data models.Metrics) error,
	data models.Metrics,
) error {
	err := fn(data)
	if err == nil || !isRetryable(err) {
		return err
	}
	retries := 100
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

func RetrySetURLData(
	fn func(nameMetric, dataMetric, typeMetric string) error,
	nameMetric, dataMetric, typeMetric string,
) error {
	err := fn(nameMetric, dataMetric, typeMetric)
	if err == nil || !isRetryable(err) {
		return err
	}
	retries := 100
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

func RetryGetAllGaugeMetrics(
	fn func() (map[string]float64, error),

) (map[string]float64, error) {
	result, err := fn()
	if err == nil || !isRetryable(err) {
		return result, err
	}

	retries := 100
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

func RetryGetAllCounterMetrics(
	fn func() (map[string]int64, error),
) (map[string]int64, error) {
	result, err := fn()
	if err == nil || !isRetryable(err) {
		return result, err
	}

	retries := 100
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

func isRetryable(err error) bool {
	if err == nil {
		return false
	}
	if r, ok := err.(RetryableError); ok {
		return r.Retryable()
	}
	return false
}
