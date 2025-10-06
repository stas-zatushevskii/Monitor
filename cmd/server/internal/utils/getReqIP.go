package utils

import (
	"net"
	"net/http"
)

func GetRequestIPAddress(r *http.Request) (string, error) {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	return host, nil
}
