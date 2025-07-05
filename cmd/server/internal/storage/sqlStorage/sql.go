package sqlStorage

import (
	"database/sql"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	storage := &PostgresStorage{db: db}
	err := storage.InitTables()
	if err != nil {
		panic(err)
	}
	return storage
}

func (ps *PostgresStorage) InitTables() error {
	_, err := ps.db.Exec(`
		CREATE TABLE IF NOT EXISTS counters (
			name TEXT PRIMARY KEY,
			value BIGINT NOT NULL
		);
		
		CREATE TABLE IF NOT EXISTS gauges (
			name TEXT PRIMARY KEY,
			value DOUBLE PRECISION NOT NULL
		);
	`)
	return err
}

func (ps *PostgresStorage) SetGauge(name string, data float64) error {
	return nil
}

func (ps *PostgresStorage) SetCounter(name string, data int64) error {
	return nil
}

func (ps *PostgresStorage) GetGauge(name string) (float64, error) {
	return 0, nil
}

func (ps *PostgresStorage) GetCounter(name string) (int64, error) {
	return 0, nil
}

func (ps *PostgresStorage) GetAllGauge() (map[string]float64, error) {
	return nil, nil
}

func (ps *PostgresStorage) GetAllCounter() (map[string]int64, error) {
	return nil, nil
}

func (ps *PostgresStorage) Ping() error  { return ps.db.Ping() }
func (ps *PostgresStorage) Close() error { return ps.db.Close() }
