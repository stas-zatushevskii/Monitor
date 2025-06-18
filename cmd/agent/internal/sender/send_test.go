package sender

import (
	"github.com/stas-zatushevskii/Monitor/cmd/agent/internal/types"
	"testing"
)

func TestCreateMetrics_Gauge(t *testing.T) {
	input := types.Gauge{
		Name: "test_gauge",
		Data: 42.42,
	}

	result, err := CreateMetrics(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != input.Name {
		t.Errorf("expected ID %s, got %s", input.Name, result.ID)
	}

	if result.MType != "gauge" {
		t.Errorf("expected MType 'gauge', got %s", result.MType)
	}

	if result.Value == nil || *result.Value != input.Data {
		t.Errorf("expected Value %f, got %v", input.Data, result.Value)
	}

	if result.Delta != nil {
		t.Errorf("expected Delta to be nil, got %v", result.Delta)
	}
}

func TestCreateMetrics_Counter(t *testing.T) {
	input := types.Counter{
		Name: "test_counter",
		Data: 123,
	}

	result, err := CreateMetrics(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != input.Name {
		t.Errorf("expected ID %s, got %s", input.Name, result.ID)
	}

	if result.MType != "counter" {
		t.Errorf("expected MType 'counter', got %s", result.MType)
	}

	if result.Delta == nil || *result.Delta != input.Data {
		t.Errorf("expected Delta %d, got %v", input.Data, result.Delta)
	}

	if result.Value != nil {
		t.Errorf("expected Value to be nil, got %v", result.Value)
	}
}
