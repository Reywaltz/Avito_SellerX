package messages

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"
)

type Message struct {
	ID        int       `json:"id"`
	Chat      int       `json:"chat"`
	Author    int       `json:"author"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

func (m *Message) PostBind(r *http.Request) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	type Alias Message

	tmp := &struct {
		Chat   *string `json:"chat"`
		Author *string `json:"author"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}

	if err := json.Unmarshal(body, &tmp); err != nil {
		return err
	}

	if tmp.Chat == nil || tmp.Author == nil {
		return errors.New("Bad request")
	}

	chatID, err := strconv.Atoi(*tmp.Chat)
	if err != nil {
		return err
	}

	authorID, err := strconv.Atoi(*tmp.Author)
	if err != nil {
		return err
	}

	m.Chat = chatID
	m.Author = authorID

	return nil
}

func (m *Message) GetBind(r *http.Request) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	type Alias Message

	tmp := &struct {
		Chat string `json:"chat"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}

	if err := json.Unmarshal(body, &tmp); err != nil {
		return err
	}

	var chatID int
	chatID, err = strconv.Atoi(tmp.Chat)
	if err != nil {
		return err
	}

	m.Chat = chatID

	return nil
}
