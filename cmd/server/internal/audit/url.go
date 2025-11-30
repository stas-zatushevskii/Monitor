package audit

import "github.com/go-resty/resty/v2"

func SendToURL(url string, data []byte) error {
	client := resty.New()
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(data).
		Post(url)
	return err
}
