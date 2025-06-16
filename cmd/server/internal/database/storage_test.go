package database

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemStorage_SetGetGauge(t *testing.T) {
	ms := NewMemStorage()
	ms.SetGauge("gauge", 1)

	res, _ := ms.GetGauge("gauge")
	assert.Equal(t, float64(1), res)

	ms.SetGauge("gauge", 10)
	resSet, _ := ms.GetGauge("gauge")
	assert.Equal(t, float64(10), resSet)
}

func TestMemStorage_SetGetCounter(t *testing.T) {
	ms := NewMemStorage()
	ms.SetCounter("counter", int64(1))

	res, _ := ms.GetCounter("counter")
	assert.Equal(t, int64(1), res)

	ms.SetCounter("counter", int64(3))
	resSet, _ := ms.GetCounter("counter")
	assert.Equal(t, int64(4), resSet)
}
