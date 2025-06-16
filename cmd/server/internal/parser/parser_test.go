package parser

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func buildRequest(body []byte) *http.Request {
	req := &http.Request{
		Body: io.NopCloser(bytes.NewReader(body)),
	}
	return req
}

func TestParseJsonData_Gauge(t *testing.T) {
	input := Metrics{
		ID:    "cpu",
		MType: "gauge",
		Value: float64Ptr(42.42),
	}

	body, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	req := buildRequest(body)

	result, err := ParseJsonData(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != input.ID {
		t.Errorf("expected ID %s, got %s", input.ID, result.ID)
	}

	if result.MType != input.MType {
		t.Errorf("expected MType %s, got %s", input.MType, result.MType)
	}

	if result.Value == nil || *result.Value != *input.Value {
		t.Errorf("expected Value %f, got %v", *input.Value, result.Value)
	}

	if result.Delta != nil {
		t.Errorf("expected Delta nil, got %v", result.Delta)
	}
}

func TestParseJsonData_Counter(t *testing.T) {
	input := Metrics{
		ID:    "requests",
		MType: "counter",
		Delta: int64Ptr(12345),
	}

	body, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	req := buildRequest(body)

	result, err := ParseJsonData(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != input.ID {
		t.Errorf("expected ID %s, got %s", input.ID, result.ID)
	}

	if result.MType != input.MType {
		t.Errorf("expected MType %s, got %s", input.MType, result.MType)
	}

	if result.Delta == nil || *result.Delta != *input.Delta {
		t.Errorf("expected Delta %d, got %v", *input.Delta, result.Delta)
	}

	if result.Value != nil {
		t.Errorf("expected Value nil, got %v", result.Value)
	}
}

func TestParseJsonData_InvalidJson(t *testing.T) {
	body := []byte(`{invalid-json}`)

	req := buildRequest(body)

	_, err := ParseJsonData(req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func int64Ptr(i int64) *int64 {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}
