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

func TestParseJSONData(t *testing.T) {
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

	result, err := ParseJSONData(req)
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

func int64Ptr(i int64) *int64 {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}
