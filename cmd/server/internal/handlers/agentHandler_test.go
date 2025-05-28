package handlers

import (
	"bytes"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/database"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAgentHandler(t *testing.T) {
	type want struct {
		statusCode  int
		contentType string
	}
	tests := []struct {
		name        string
		method      string
		contentType string
		url         string
		want        want
	}{
		{
			name:        "valid gauge metric",
			method:      http.MethodPost,
			contentType: "text/plain",
			url:         "/update/gauge/temperature/23.5",
			want: want{
				statusCode:  http.StatusOK,
				contentType: "",
			},
		},
		{
			name:        "invalid method",
			method:      http.MethodGet,
			contentType: "text/plain",
			url:         "/update/gauge/temperature/23.5",
			want: want{
				statusCode:  http.StatusMethodNotAllowed,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "invalid content type",
			method:      http.MethodPost,
			contentType: "application/json",
			url:         "/update/gauge/temperature/23.5",
			want: want{
				statusCode:  http.StatusUnsupportedMediaType,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, bytes.NewBuffer(nil))
			req.Header.Set("Content-Type", tt.contentType)

			rec := httptest.NewRecorder()
			handler := UpdateAgentHandler(database.NewMemStorage())
			handler.ServeHTTP(rec, req)

			resp := rec.Result()
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)

			if resp.StatusCode != tt.want.statusCode {
				t.Errorf("[%s] expected status %d, got %d", tt.name, tt.want.statusCode, resp.StatusCode)
			}
			if tt.want.contentType != "" && resp.Header.Get("Content-Type") != tt.want.contentType {
				t.Errorf("[%s] expected content-type %s, got %s", tt.name, tt.want.contentType, resp.Header.Get("Content-Type"))
			}
			if resp.StatusCode != http.StatusOK && len(body) == 0 {
				t.Errorf("[%s] expected error message in response body, got none", tt.name)
			}
		})
	}
}
