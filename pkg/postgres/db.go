package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

type DB struct {
	conn *pgx.Conn
}

func (db *DB) Conn() *pgx.Conn {
	return db.conn
}

func NewDB() (*DB, error) {
	connstring := os.Getenv("CONN_DB")
	if connstring == "" {
		connstring = "postgres://avito_user:avito_pass@localhost:5433/avito"
	}
	conn, err := pgx.Connect(context.Background(), connstring)
	if err != nil {
		return nil, fmt.Errorf("Can't init connection to db: %w", err)
	}

	return &DB{conn: conn}, nil
}
