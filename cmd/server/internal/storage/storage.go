package storage

// in memory storage

type CounterStorage interface {
	SetCounter(name string, data int64) error
	GetCounter(name string) (int64, error)
	GetAllCounter() (map[string]int64, error)
}
type GaugeStorage interface {
	SetGauge(name string, data float64) error
	GetGauge(name string) (float64, error)
	GetAllGauge() (map[string]float64, error)
}

// abstract storage

type Storage interface {
	CounterStorage
	GaugeStorage
	Ping() error
	Close() error
}
