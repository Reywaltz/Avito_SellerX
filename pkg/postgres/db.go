package postgres

import (
	"context"
	"fmt"

	"github.com/Reywaltz/avito_backend/configs"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
}

func (db *DB) Pool() *pgxpool.Pool {
	return db.pool
}

func NewDB(cfg configs.Config) (*DB, error) {
	conn, err := pgxpool.Connect(context.Background(), cfg.ConnString)
	if err != nil {
		return nil, fmt.Errorf("Can't init connection to db: %w", err)
	}

	return &DB{pool: conn}, nil
}
