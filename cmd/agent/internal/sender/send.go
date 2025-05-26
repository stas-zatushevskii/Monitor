package sender

import (
	"github.com/stas-zatushevskii/Monitor/cmd/agent/internal/types"
	"log"
	"net/http"
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
	req, err := http.NewRequest(http.MethodPost, newURL, nil)
	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Set("Content-Type", "text/plain")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("failed to send request: %v", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("failed to send request: %v, [%s]", resp.Status, url)
	}
}
