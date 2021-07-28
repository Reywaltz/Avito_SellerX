package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Reywaltz/avito_backend/internal/message"
	"github.com/Reywaltz/avito_backend/internal/models/chats"
	"github.com/Reywaltz/avito_backend/internal/models/users"
	log "github.com/Reywaltz/avito_backend/pkg/log"
	"github.com/Reywaltz/avito_backend/pkg/postgres"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
)

type ChatRepository interface {
	Create(chat chats.Chat) (int, error)
	GetChats(user users.User) ([]chats.Chat, error)
}

type ChatHandlers struct {
	Log      log.Logger
	ChatRepo ChatRepository
	UserRepo UserRepository
}

func NewChatHandlers(logger log.Logger, chatRepo ChatRepository, userRepo UserRepository) *ChatHandlers {
	return &ChatHandlers{
		Log:      logger,
		ChatRepo: chatRepo,
		UserRepo: userRepo,
	}
}

func (q *ChatHandlers) Create(w http.ResponseWriter, r *http.Request) {
	var chat chats.Chat

	err := chat.Bind(r)
	if err != nil {
		q.Log.Errorf("Can't bind chat: {%s}", err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	for _, tmpID := range chat.Users {
		userID, _ := strconv.Atoi(tmpID)
		tmpUser := users.User{ID: userID}
		_, err = q.UserRepo.GetOne(tmpUser)
		if err != nil {
			q.Log.Errorf("Not exsisting user by id: {%d}", tmpUser.ID)
			w.WriteHeader(http.StatusBadRequest)

			return
		}
	}

	chatID, err := q.ChatRepo.Create(chat)
	if err != nil {
		if errors.Is(err, postgres.DuplicateError) {
			w.WriteHeader(http.StatusBadRequest)
			q.Log.Errorf("Can't create new chat: %s", err)
			return
		}
	}
	q.Log.Infof("Created chat with ID: {%d}", chatID)
	message.MakeResponse(w, chatID, http.StatusCreated)
}

func (q *ChatHandlers) GetChats(w http.ResponseWriter, r *http.Request) {
	var user users.User

	if err := user.GetBind(r); err != nil {
		q.Log.Errorf("Can't bind json: %s", err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	_, err := q.UserRepo.GetOne(user)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			q.Log.Errorf("User id=%d doesn't exist", user.ID)
			w.WriteHeader(http.StatusNotFound)

			return
		}
		q.Log.Errorf("Can't get data from db: %s", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	res, err := q.ChatRepo.GetChats(user)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			q.Log.Infof("User %d don't have any messages", user.ID)
			w.WriteHeader(http.StatusNotFound)

			return
		}
		q.Log.Errorf("Can't get chats: %s", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	out, err := json.Marshal(res)
	if err != nil {
		q.Log.Errorf("Can't marshall: %s", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func (q *ChatHandlers) Route(router *mux.Router) {
	s := router.PathPrefix("/chats").Subrouter()
	s.HandleFunc("/add", q.Create).Methods("POST")
	s.HandleFunc("/get", q.GetChats).Methods("POST")
}
