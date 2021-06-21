package message_repo

import (
	"context"
	"time"

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

	res := r.db.Pool().QueryRow(context.Background(), createMessage,
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

	res, err := r.db.Pool().Query(context.Background(), GetMessages, message.Chat)
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
	Chatmessagesquery = `select ` + selectMessageFields + `from messages where chat = $1 ORDER BY created_at DESC`
)

func (r *MessageRepo) GetChatMessages(message messages.Message) ([]messages.Message, error) {
	out := make([]messages.Message, 0)

	res, err := r.db.Pool().Query(context.Background(), Chatmessagesquery, message.Chat)
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
