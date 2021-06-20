package messages

import "time"

type Message struct {
	ID        int       `json:"id"`
	Chat      int       `json:"chat"`
	Author    int       `json:"author"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}
