package sqlstorage

import (
	"context"
	"database/sql"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/models"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

func (ps *PostgresStorage) Bootstrap(ctx context.Context) error {
	tx, err := ps.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS counters (
			name TEXT PRIMARY KEY,
			value BIGINT NOT NULL
		)`); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS gauges (
			name TEXT PRIMARY KEY,
			value DOUBLE PRECISION NOT NULL
		);`); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
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
			ON CONFLICT (name) DO UPDATE 
			SET value = EXCLUDED.value
			WHERE EXCLUDED.value > counters.value
	`, name, data)
	return err
}

func (ps *PostgresStorage) GetGauge(name string) (float64, error) {
	var value float64
	err := ps.db.QueryRow(`SELECT value FROM gauges WHERE name = $1`, name).Scan(&value)
	if err != nil {
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

func (ps *PostgresStorage) SetMultipleGauge(ctx context.Context, metrics []models.Metrics) error {
	tx, err := ps.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO gauges (name, value)
		VALUES ($1, $2)
		ON CONFLICT (name) DO UPDATE SET value = $2
	`)
	if err != nil {
		return err
	}

	for _, v := range metrics {
		if v.Value == nil {
			continue
		}
		_, err := stmt.ExecContext(ctx, v.ID, *v.Value)
		if err != nil {
			return err
		}
	}
	err = stmt.Close()
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (ps *PostgresStorage) SetMultipleCounter(ctx context.Context, metrics []models.Metrics) error {
	tx, err := ps.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
			INSERT INTO counters (name, value)
			VALUES ($1, $2)
			ON CONFLICT (name) DO UPDATE 
			SET value = EXCLUDED.value + counters.value
	`)
	if err != nil {
		return err
	}

	for _, v := range metrics {
		if v.Delta == nil {
			continue
		}
		_, err := stmt.ExecContext(ctx, v.ID, *v.Delta)
		if err != nil {
			return err
		}
	}
	err = stmt.Close()
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (ps *PostgresStorage) Ping() error  { return ps.db.Ping() }
func (ps *PostgresStorage) Close() error { return ps.db.Close() }
