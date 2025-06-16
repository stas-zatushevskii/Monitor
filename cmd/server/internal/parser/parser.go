package parser

import (
	"encoding/json"
	"fmt"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/constants"
	"net/http"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение мxетрики в случае передачи gauge
}

func ParseJsonData(r *http.Request) (Metrics, error) {
	var data Metrics
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return Metrics{}, fmt.Errorf(constants.ErrorParseJson)
	}
	defer r.Body.Close()
	return data, nil
}
