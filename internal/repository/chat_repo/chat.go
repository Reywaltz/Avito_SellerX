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

const (
	GetUsersChat = `select ` + selectUserFields + ` from users_chats 
	inner join chats on users_chats.chat_id = chats.id where user_id = $1 ORDER BY created_at DESC`
)

func (r *ChatRepo) GetChats(userID int) ([]chats.Chat, error) {
	var out []chats.Chat

	res, err := r.db.Pool().Query(context.Background(), GetUsersChat, userID)
	if err != nil {
		return nil, err
	}

	for res.Next() {
		var tmp chats.Chat

		if err = res.Scan(&tmp.ID, &tmp.Name, &tmp.CreatedAt); err != nil {
			return nil, err
		}

		out = append(out, tmp)
	}

	return out, nil
}
