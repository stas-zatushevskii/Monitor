package sender

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/stas-zatushevskii/Monitor/cmd/agent/internal/types"
	"strconv"
)

func CreatePath[metricData types.Gauge | types.Counter](m metricData, url string) string {
	// http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
	var metricValue string
	switch v := any(m).(type) {
	case types.Counter:
		metricValue = strconv.FormatInt(v.Data, 10)
		url += "/update/counter/" + v.Name + "/" + metricValue
	case types.Gauge:
		metricValue = strconv.FormatFloat(v.Data, 'f', -1, 64)
		url += "/update/gauge/" + v.Name + "/" + metricValue
	default:
		return ""
	}
	return url
}

func SendData[metricData types.Gauge | types.Counter](m metricData, url string) {
	var newURL = CreatePath(m, url)
	client := resty.New()
	_, err := client.R().
		SetHeader("Content-Type", "text/plainn").
		Post(newURL)

	if err != nil {
		fmt.Println(err)
		return
	}
}
