package message

import (
	"context"
	"time"

	"github.com/Reywaltz/avito_backend/internal/models/chats"
	"github.com/Reywaltz/avito_backend/internal/models/messages"
	"github.com/Reywaltz/avito_backend/pkg/postgres"
)

const (
	messageFields       = `chat, author, text, created_at`
	selectMessageFields = `id, ` + messageFields
)

type MessageRepo struct {
	db *postgres.DB
}

func NewMessageRepository(db *postgres.DB) *MessageRepo {
	return &MessageRepo{
		db: db,
	}
}

const (
	createMessage = `INSERT INTO messages ( ` + messageFields + `) values ($1, $2, $3, $4) RETURNING id`
)

func (r *MessageRepo) Create(message messages.Message) (int, error) {
	var messageID int

	message.CreatedAt = time.Now()

	res := r.db.Conn().QueryRow(context.Background(), createMessage,
		message.Chat,
		message.Author,
		message.Text,
		message.CreatedAt)

	if err := res.Scan(&messageID); err != nil {
		return 0, err
	}

	return messageID, nil
}

const (
	GetMessages = `SELECT ` + selectMessageFields + ` FROM messages WHERE chat=$1 ORDER BY created_at DESC`
)

func (r *MessageRepo) GetMessages(message messages.Message) ([]messages.Message, error) {
	out := make([]messages.Message, 0)

	res, err := r.db.Conn().Query(context.Background(), GetMessages, message.Chat)
	if err != nil {
		return nil, err
	}

	for res.Next() {
		var curMessage messages.Message
		if err := res.Scan(&curMessage.ID, &curMessage.Chat, &curMessage.Author, &curMessage.Text, &curMessage.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, curMessage)
	}

	return out, nil
}

const (
	GetOne = `SELECT * FROM chats WHERE id = $1`
)

func (r *MessageRepo) GetOne(message messages.Message) (chats.Chat, error) {
	var chat chats.Chat

	res := r.db.Conn().QueryRow(context.Background(), GetOne, message.Chat)
	if err := res.Scan(&chat.ID, &chat.Name, &chat.CreatedAt); err != nil {
		return chat, err
	}

	return chat, nil
}

const (
	GetUserInChat = "select chat_id, user_id from users_chats where user_id = $1 and chat_id=$2"
)

type UserChat struct {
	ChatID int
	UserID int
}

func (r *MessageRepo) CheckUser(message messages.Message) bool {
	res := r.db.Conn().QueryRow(context.Background(), GetUserInChat, message.Author, message.Chat)

	var tmp UserChat
	if err := res.Scan(tmp); err != nil {
		return false
	}

	return true
}
