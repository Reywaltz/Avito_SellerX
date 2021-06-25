package user_repo

import (
	"context"

	"github.com/Reywaltz/avito_backend/internal/models/users"
	"github.com/Reywaltz/avito_backend/pkg/postgres"
)

const (
	userFields       = `username, created_at`
	selectUserFields = `id, ` + userFields
)

type UserRepo struct {
	db *postgres.DB
}

func NewUserRepository(db *postgres.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

const createUser = `INSERT INTO users ( ` + userFields + `) values ($1, $2) RETURNING id`

func (r *UserRepo) Create(user users.User) (int, error) {
	var createdID int

	row := r.db.Conn().QueryRow(context.Background(), createUser, user.Username, user.CreatedAt)

	if err := row.Scan(&createdID); err != nil {
		if postgres.IsDuplicated(err) {
			return 0, postgres.DuplicateError
		}

		return 0, err
	}

	return createdID, nil
}

const (
	selectUsers    = `SELECT ` + selectUserFields + ` FROM users`
	selectUserByID = `SELECT ` + selectUserFields + ` FROM users WHERE id=$1`
)

func (r *UserRepo) GetAll() ([]users.User, error) {
	res, err := r.db.Conn().Query(context.Background(), selectUsers)
	if err != nil {

		return nil, err
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

func (r *UserRepo) GetOne(user users.User) (users.User, error) {
	res := r.db.Conn().QueryRow(context.Background(), selectUserByID, user.ID)

	var out users.User
	err := res.Scan(&out.ID, &out.Username, &out.CreatedAt)
	if err != nil {
		return out, err
	}

	return out, nil
}
