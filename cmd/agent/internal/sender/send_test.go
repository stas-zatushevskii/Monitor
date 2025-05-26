package sender

import (
	"github.com/stas-zatushevskii/Monitor/cmd/agent/internal/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreatePathGauge(t *testing.T) {
	type want struct {
		path string
	}
	tests := []struct {
		name string
		data types.Gauge
		want want
	}{
		{
			name: "TestGaugePath",
			data: types.Gauge{Name: "b", Data: 1.1},
			want: want{
				path: "/update/gauge/b/1.1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want.path, CreatePath(tt.data, ""))
		})
	}
}

func TestCreatePathCounter(t *testing.T) {
	type want struct {
		path string
	}
	tests := []struct {
		name string
		data types.Counter
		want want
	}{
		{
			name: "TestGaugePath",
			data: types.Counter{Name: "b", Data: 1},
			want: want{
				path: "/update/counter/b/1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want.path, CreatePath(tt.data, ""))
		})
	}
}
