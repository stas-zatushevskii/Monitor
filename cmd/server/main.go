package main

import (
	"net/http"
	"strconv"
)

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

func (ms *MemStorage) setGauge(name string, data float64) {
	ms.gauge[name] = data
}

func (ms *MemStorage) setCounter(name string, data int64) {
	ms.counter[name] = data
}

func (ms *MemStorage) getGauge(name string) float64 {
	return ms.gauge[name]
}

func (ms *MemStorage) getCounter(name string) int64 {
	return ms.counter[name]
}

func agentHandler(w http.ResponseWriter, r *http.Request) {
	storage := NewMemStorage()
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "Content-Type not allowed", http.StatusUnsupportedMediaType)
		return
	}
	if err := r.ParseForm(); err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	nameMetric := r.FormValue("name")
	if nameMetric == "" {
		http.Error(w, "Missing metric name", http.StatusNotFound)
		return
	}

	dataMetric := r.FormValue("data")
	typeMetric := r.FormValue("type")

	switch typeMetric {
	case "gauge":
		parsedData, err := strconv.ParseFloat(dataMetric, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		storage.setGauge(nameMetric, parsedData)
		w.WriteHeader(http.StatusOK)
		return
	case "counter":
		parsedData, err := strconv.ParseInt(dataMetric, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		storage.setCounter(nameMetric, parsedData)
		w.WriteHeader(http.StatusOK)
		return
	default:
		http.Error(w, "Unsupported type", http.StatusBadRequest)
		return
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", agentHandler)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
