package postgres

import (
	"errors"

	"github.com/jackc/pgconn"
)

const (
	duplicateErrorCode = "23505"
)

var DuplicateError = errors.New("value in database already exists")

func IsDuplicated(err error) bool {
	var pgxerr *pgconn.PgError

	if errors.As(err, &pgxerr) {
		if pgxerr.Code == duplicateErrorCode {
			return true
		}
	}

	return false
}
