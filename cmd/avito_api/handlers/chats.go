package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Reywaltz/avito_backend/internal/models/chats"
	"github.com/Reywaltz/avito_backend/internal/models/users"
	log "github.com/Reywaltz/avito_backend/pkg/log"
	"github.com/Reywaltz/avito_backend/pkg/postgres"
	"github.com/gorilla/mux"
)

type ChatRepository interface {
	Create(chat chats.Chat) (int, error)
	GetChats(userID int) ([]chats.Chat, error)
}

type ChatHandlers struct {
	Log           log.Logger
	ChatRepo      ChatRepository
	UsersChatRepo UserRepository
}

func NewChatHandlers(logger log.Logger, chatRepo ChatRepository, userRepo UserRepository) *ChatHandlers {
	return &ChatHandlers{
		Log:           logger,
		ChatRepo:      chatRepo,
		UsersChatRepo: userRepo,
	}
}

func (q *ChatHandlers) Create(writer http.ResponseWriter, request *http.Request) {
	var chat chats.Chat

	chat.Bind(request)

	for _, item := range chat.Users {
		tmp, err := strconv.Atoi(item)
		if err != nil {
			q.Log.Errorf("Can't parse user id: {%s}", item)
			writer.WriteHeader(http.StatusBadRequest)

			return
		}

		_, err = q.UsersChatRepo.GetOne(tmp)
		if err != nil {
			q.Log.Errorf("Not exsisting user by id: {%d}", tmp)
			writer.WriteHeader(http.StatusBadRequest)

			return
		}
	}

	chatID, err := q.ChatRepo.Create(chat)
	if err != nil {
		if errors.Is(err, postgres.DuplicateError) {
			writer.WriteHeader(http.StatusBadRequest)
			q.Log.Errorf("Can't create new chat: %s", err)
			return
		}
	}

	q.Log.Infof("Created chat with ID: {%d}", chatID)

	writer.WriteHeader(http.StatusCreated)

}

func (q *ChatHandlers) GetChats(writer http.ResponseWriter, request *http.Request) {
	var user users.User

	if err := user.GetBind(request); err != nil {
		q.Log.Errorf("Can't bind json: %s", err)
	}

	res, err := q.ChatRepo.GetChats(user.ID)
	if err != nil {
		q.Log.Errorf("Can't get chats: %s", err)

		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	out, err := json.Marshal(res)
	if err != nil {
		q.Log.Errorf("Can't marshall: %s", err)

		writer.WriteHeader(http.StatusInternalServerError)
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(out)

	return
}

func (q *ChatHandlers) Route(router *mux.Router) {
	s := router.PathPrefix("/chats").Subrouter()
	s.HandleFunc("/add", q.Create).Methods("POST")
	s.HandleFunc("/get", q.GetChats).Methods("POST")
}
