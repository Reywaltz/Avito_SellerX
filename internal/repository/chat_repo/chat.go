package chat_repo

import (
	"context"
	"time"

	"github.com/Reywaltz/avito_backend/internal/models/chats"
	log "github.com/Reywaltz/avito_backend/pkg/log"
	"github.com/Reywaltz/avito_backend/pkg/postgres"
)

const (
	chatFields       = `name, created_at`
	selectUserFields = `id, ` + chatFields
)

type ChatRepo struct {
	db     *postgres.DB
	logger log.Logger
}

func NewChatRepository(db *postgres.DB, logger log.Logger) *ChatRepo {
	return &ChatRepo{
		db:     db,
		logger: logger,
	}
}

const (
	createChat      = `INSERT INTO chats ( ` + chatFields + `) values ($1, $2) RETURNING id`
	createdChatUser = `INSERT INTO users_chats (chat_id, user_id) values ($1, $2)`
)

func (r *ChatRepo) Create(chat chats.Chat) (int, error) {
	var chatID int

	chat.CreatedAt = time.Now()
	tx, err := r.db.Pool().Begin(context.Background())
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(context.Background())

	res := tx.QueryRow(context.Background(), createChat, chat.Name, chat.CreatedAt)
	if err = res.Scan(&chatID); err != nil {
		if postgres.IsDuplicated(err) {
			return 0, postgres.DuplicateError
		}

		return 0, err
	}
	for _, userID := range chat.Users {
		_, err = tx.Exec(context.Background(), createdChatUser, chatID, userID)
		if err != nil {
			return 0, err
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return 0, err
	}

	return chatID, nil
}
