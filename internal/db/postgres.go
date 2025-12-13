package db

import (
	"context"
	"database/sql"

	// Postgres driver
	_ "github.com/lib/pq" // Driver registration
	"github.com/singh-anurag-7991/shield/internal/models"
)

type PostgresConfigStore struct {
	db *sql.DB
}

func NewPostgresConfigStore(dsn string) (*PostgresConfigStore, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresConfigStore{db: db}, nil
}

func (p *PostgresConfigStore) LoadConfigs(ctx context.Context) ([]models.LimiterConfig, error) {
	rows, err := p.db.QueryContext(ctx, `SELECT name, algorithm, capacity, rate, "window" FROM limiter_configs`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []models.LimiterConfig
	for rows.Next() {
		var cfg models.LimiterConfig
		var window sql.NullString
		if err := rows.Scan(&cfg.Name, &cfg.Algorithm, &cfg.Capacity, &cfg.Rate, &window); err != nil {
			return nil, err
		}
		if window.Valid {
			cfg.Window = window.String
		} // If NULL, Window remains empty string
		configs = append(configs, cfg)
	}
	return configs, rows.Err()
}

func (p *PostgresConfigStore) InitTable(ctx context.Context) error {
	_, err := p.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS limiter_configs (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) UNIQUE NOT NULL,
			algorithm VARCHAR(50) NOT NULL,
			capacity BIGINT NOT NULL,
			rate BIGINT,
			"window" VARCHAR(50)
		);
		INSERT INTO limiter_configs (name, algorithm, capacity, rate, "window") VALUES
			('global', 'token', 10, 10, NULL),
			('burst', 'leaky', 5, 2, NULL)
		ON CONFLICT (name) DO NOTHING;
	`)
	return err
}
