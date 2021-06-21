package chats

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type Chat struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Users     []string  `json:"users,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *Chat) Bind(r *http.Request) error {
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
