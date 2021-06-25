package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	connstring = "postgres://avito_user:avito_pass@localhost:5433/avito"
)

type DB struct {
	pool *pgxpool.Pool
}

func (db *DB) Pool() *pgxpool.Pool {
	return db.pool
}

func NewDB() (*DB, error) {
	conn, err := pgxpool.Connect(context.Background(), connstring)
	if err != nil {
		return nil, fmt.Errorf("Can't init connection to db: %w", err)
	}

	return &DB{pool: conn}, nil
}
