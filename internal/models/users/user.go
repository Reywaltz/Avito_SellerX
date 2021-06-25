package users

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
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

func (u *User) GetBind(r *http.Request) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	type Alias User

	tmp := &struct {
		User *string `json:"user"`
		*Alias
	}{
		Alias: (*Alias)(u),
	}

	if err := json.Unmarshal(body, &tmp); err != nil {
		return err
	}

	if tmp.User == nil {
		return errors.New("UserID is empty")
	}

	userID, err := strconv.Atoi(*tmp.User)
	if err != nil {
		return err
	}

	u.ID = userID

	return nil
}
