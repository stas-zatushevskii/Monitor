package utils

import (
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/models"
)

func GetMetricName(data ...models.Metrics) []string {
	result := make([]string, len(data))
	for _, m := range data {
		if m.ID != "" {
			result = append(result, m.MType)
		}
		if m.MType != "" {
			result = append(result, m.MType)
		}
	}
	return result
}
