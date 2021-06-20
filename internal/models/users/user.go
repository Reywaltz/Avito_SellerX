package users

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *User) Bind(r *http.Request) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &u)
	if err != nil {
		return err
	}

	return nil
}
