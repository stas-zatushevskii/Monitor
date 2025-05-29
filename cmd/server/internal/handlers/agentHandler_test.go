package handlers

import (
	"bytes"
	"github.com/go-chi/chi/v5"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/database"
	"net/http"
	"net/http/httptest"
	"testing"
)

func RouterForTest(storage *database.MemStorage) *chi.Mux {
	router := chi.NewRouter()

	router.Get("/", GetAllAgentHandlers(storage))
	router.Post("/update/{type}/{name}/{data}", UpdateAgentHandler(storage))
	router.Get("/value/{type}/{name}", ValueAgentHandler(storage))
	return router
}

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
				statusCode: http.StatusOK,
			},
		},
		{
			name:        "invalid method",
			method:      http.MethodGet,
			contentType: "text/plain",
			url:         "/update/gauge/temperature/23.5",
			want: want{
				statusCode: http.StatusMethodNotAllowed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, bytes.NewBuffer(nil))
			req.Header.Set("Content-Type", tt.contentType)

			rec := httptest.NewRecorder()

			storage := database.NewMemStorage()
			router := RouterForTest(storage)
			router.ServeHTTP(rec, req)

			resp := rec.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tt.want.statusCode {
				t.Errorf("[%s] expected status %d, got %d", tt.name, tt.want.statusCode, resp.StatusCode)
			}
			if tt.want.contentType != "" && resp.Header.Get("Content-Type") != tt.want.contentType {
				t.Errorf("[%s] expected content-type %s, got %s", tt.name, tt.want.contentType, resp.Header.Get("Content-Type"))
			}
		})
	}
}
