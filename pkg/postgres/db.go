package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
}

type Config struct {
	ConnString string `json:"connstring"`
}

func InitConfig() (Config, error) {
	file, err := os.ReadFile("cmd/avito_api/configs/cfg.json")
	if err != nil {
		return Config{}, fmt.Errorf("Can't open file: %w", err)
	}
	var cfg Config
	if err = json.Unmarshal(file, &cfg); err != nil {
		return Config{}, fmt.Errorf("Can't unmarshall json file: %w", err)
	}

	return cfg, nil
}

func (db *DB) Pool() *pgxpool.Pool {
	return db.pool
}

func NewDB(cfg Config) (*DB, error) {
	conn, err := pgxpool.Connect(context.Background(), cfg.ConnString)
	if err != nil {
		return nil, fmt.Errorf("Can't init connection to db: %w", err)
	}

	return &DB{pool: conn}, nil
}
