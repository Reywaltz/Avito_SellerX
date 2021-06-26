package chats

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"
)

type Chat struct {
	ID        int       `json:"id"`
	Name      *string   `json:"name"`
	Users     []string  `json:"users,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func (c *Chat) Bind(r *http.Request) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &c); err != nil {
		return err
	}

	if c.Name == nil {
		return errors.New("Name must be provided")
	}

	if c.Users == nil {
		return errors.New("Users must be provided")
	}

	for _, userID := range c.Users {
		_, err = strconv.Atoi(userID)
		if err != nil {
			return err
		}
	}

	return nil
}
