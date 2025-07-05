package sqlstorage

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
	_, err := ps.db.Exec(`
		INSERT INTO gauges (name, value)
		VALUES ($1, $2)
		ON CONFLICT (name) DO UPDATE SET value = $2
	`, name, data)
	return err
}

func (ps *PostgresStorage) SetCounter(name string, data int64) error {
	_, err := ps.db.Exec(`
		INSERT INTO counters (name, value)
		VALUES ($1, $2)
		ON CONFLICT (name) DO UPDATE SET value = counters.value + EXCLUDED.value
	`, name, data)
	return err
}

func (ps *PostgresStorage) GetGauge(name string) (float64, error) {
	var value float64
	err := ps.db.QueryRow(`SELECT value FROM gauges WHERE name = $1`, name).Scan(&value)
	if err == sql.ErrNoRows {
		return 0, nil // по аналогии с in-memory можно вернуть 0 без ошибки
	}
	return value, err
}

func (ps *PostgresStorage) GetCounter(name string) (int64, error) {
	var value int64
	err := ps.db.QueryRow(`SELECT value FROM counters WHERE name = $1`, name).Scan(&value)
	if err != nil {
		return 0, nil
	}
	return value, err
}

func (ps *PostgresStorage) GetAllGauge() (map[string]float64, error) {
	rows, err := ps.db.Query(`SELECT name, value FROM gauges`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]float64)
	for rows.Next() {
		var name string
		var value float64
		if err := rows.Scan(&name, &value); err != nil {
			return nil, err
		}
		result[name] = value
	}
	return result, rows.Err()
}

func (ps *PostgresStorage) GetAllCounter() (map[string]int64, error) {
	rows, err := ps.db.Query(`SELECT name, value FROM counters`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]int64)
	for rows.Next() {
		var name string
		var value int64
		if err := rows.Scan(&name, &value); err != nil {
			return nil, err
		}
		result[name] = value
	}
	return result, rows.Err()
}
