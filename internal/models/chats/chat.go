package chats

import (
	"time"
)

type Chat struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Users     []string  `json:"users"`
	CreatedAt time.Time `json:"created_at"`
}
