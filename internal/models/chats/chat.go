package chats

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"
)

type Chat struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Users     []string  `json:"users,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func (c *Chat) Bind(r *http.Request) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	type Alias Chat

	tmp := &struct {
		Chat *string `json:"chat"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(body, &tmp); err != nil {
		return err
	}

	var chatID int
	chatID, err = strconv.Atoi(*tmp.Chat)
	if err != nil {
		return err
	}

	c.ID = chatID

	return nil
}
