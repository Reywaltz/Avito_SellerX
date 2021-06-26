package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

const (
	connstring = "postgres://avito_user:avito_pass@db:5432/avito"
)

type DB struct {
	conn *pgx.Conn
}

func (db *DB) Conn() *pgx.Conn {
	return db.conn
}

func NewDB() (*DB, error) {
	conn, err := pgx.Connect(context.Background(), connstring)
	if err != nil {
		return nil, fmt.Errorf("Can't init connection to db: %w", err)
	}

	return &DB{conn: conn}, nil
}
