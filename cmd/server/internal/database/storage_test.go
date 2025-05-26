package database

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestMemStorage_SetGetGauge(t *testing.T) {
	ms := NewMemStorage()
	ms.SetGauge("gauge", 1)

	ms.GetGauge("gauge")
	assert.Equal(t, float64(1), ms.GetGauge("gauge"))

	ms.SetGauge("gauge", 10)
	assert.Equal(t, float64(10), ms.GetGauge("gauge"))
}

func TestMemStorage_SetGetCounter(t *testing.T) {
	ms := NewMemStorage()
	ms.SetCounter("counter", int64(1))

	ms.GetCounter("counter")
	assert.Equal(t, int64(1), ms.GetCounter("counter"))

	ms.SetCounter("counter", int64(3))
	assert.Equal(t, int64(4), ms.GetCounter("counter"))
}

func TestParseData(t *testing.T) {
	type want struct {
		nameMetric string
		dataMetric string
		typeMetric string
		err        error
	}
	tests := []struct {
		name string
		url  string
		want want
	}{
		{
			name: "success",
			url:  "http://localhost/update/counter/a/1",
			want: want{
				nameMetric: "a",
				dataMetric: "1",
				typeMetric: "counter",
				err:        nil,
			},
		},
		{
			name: "wrong url",
			url:  "http://localhost/1",
			want: want{
				nameMetric: "",
				dataMetric: "",
				typeMetric: "",
				err:        errors.New("invalid URL path: /1"),
			},
		},
		{
			name: "missing nameMetric",
			url:  "http://localhost/update/counter//1",
			want: want{
				nameMetric: "",
				dataMetric: "",
				typeMetric: "",
				err:        errors.New("missing metric name in URL"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, _ := url.Parse(tt.url)
			gotName, gotData, gotType, gotErr := ParseData(u.Path)

			assert.Equal(t, tt.want.nameMetric, gotName)
			assert.Equal(t, tt.want.dataMetric, gotData)
			assert.Equal(t, tt.want.typeMetric, gotType)
			if tt.want.err != nil {
				assert.EqualError(t, gotErr, tt.want.err.Error())
			} else {
				assert.NoError(t, gotErr)
			}
		})
	}
}
