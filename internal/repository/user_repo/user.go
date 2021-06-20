package user_repo

import (
	"context"
	"fmt"

	"github.com/Reywaltz/avito_backend/internal/models/users"
	log "github.com/Reywaltz/avito_backend/pkg/log"
	"github.com/Reywaltz/avito_backend/pkg/postgres"
)

const (
	userFields       = `username, created_at`
	selectUserFields = `id, ` + userFields
)

type UserRepo struct {
	db     *postgres.DB
	logger log.Logger
}

func NewUserRepository(db *postgres.DB, logger log.Logger) *UserRepo {
	return &UserRepo{
		db:     db,
		logger: logger,
	}
}

const createUser = `INSERT INTO users ( ` + userFields + `) values ($1, $2) RETURNING id`

func (r *UserRepo) Create(user users.User) (int, error) {
	var createdID int

	row := r.db.Pool().QueryRow(context.Background(), createUser, user.Username, user.CreatedAt)

	if err := row.Scan(&createdID); err != nil {
		if postgres.IsDuplicated(err) {
			return -1, postgres.DuplicateError
		}

		return -1, err
	}

	return createdID, nil
}

const (
	selectUsers    = `SELECT ` + selectUserFields + ` FROM users`
	selectUserByID = `SELECT ` + selectUserFields + ` FROM users WHERE id=$1`
)

func (r *UserRepo) GetAll() ([]users.User, error) {
	res, err := r.db.Pool().Query(context.Background(), selectUsers)
	if err != nil {
		r.logger.Errorf("Can't init query to database: %s", err.Error())

		return nil, fmt.Errorf("Can't init query to database: %s", err)
	}

	out := make([]users.User, 0)
	for res.Next() {
		var users users.User
		err := res.Scan(&users.ID, &users.Username, &users.CreatedAt)
		if err != nil {
			return nil, err
		}
		out = append(out, users)
	}

	return out, nil
}

func (r *UserRepo) GetOne(id int) (users.User, error) {
	res := r.db.Pool().QueryRow(context.Background(), selectUserByID, id)

	var user users.User
	err := res.Scan(&user.ID, &user.Username, &user.CreatedAt)
	if err != nil {
		return user, err
	}

	return user, nil
}
