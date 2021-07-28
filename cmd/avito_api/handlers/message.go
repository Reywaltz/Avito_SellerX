package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Reywaltz/avito_backend/internal/message"
	"github.com/Reywaltz/avito_backend/internal/models/chats"
	"github.com/Reywaltz/avito_backend/internal/models/messages"
	log "github.com/Reywaltz/avito_backend/pkg/log"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
)

type MessageRepository interface {
	Create(message messages.Message) (int, error)
	GetMessages(message messages.Message) ([]messages.Message, error)
	GetOne(message messages.Message) (chats.Chat, error)
	CheckUser(message messages.Message) bool
}

type MessageHandlers struct {
	Log         log.Logger
	MessageRepo MessageRepository
}

func NewMessageHandlers(logger log.Logger, messageRepo MessageRepository) *MessageHandlers {
	return &MessageHandlers{
		Log:         logger,
		MessageRepo: messageRepo,
	}
}

func (q *MessageHandlers) Create(w http.ResponseWriter, r *http.Request) {
	var newMessage messages.Message

	if err := newMessage.PostBind(r); err != nil {
		q.Log.Errorf("Can't insert message: %s", err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if q.MessageRepo.CheckUser(newMessage) {
		messageID, err := q.MessageRepo.Create(newMessage)
		if err != nil {
			q.Log.Errorf("Can't create new message: %s", err)
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		q.Log.Infof("New message id: %d", messageID)
		message.MakeResponse(w, messageID, http.StatusCreated)

		return
	}

	q.Log.Errorf("User %d is not present in chat %d", newMessage.Author, newMessage.Chat)
	rawJSON := []byte(`{"error": "user not present in chat"}`)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(rawJSON)

	return
}

func (q *MessageHandlers) GetMessages(w http.ResponseWriter, r *http.Request) {
	var message messages.Message

	if err := message.GetBind(r); err != nil {
		q.Log.Errorf("Error bind: %s", err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	_, err := q.MessageRepo.GetOne(message)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			q.Log.Errorf("Chat is not found: %s", err)
			w.WriteHeader(http.StatusNotFound)

			return
		} else {
			q.Log.Errorf("Can't get data from db: %s", err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}
	}

	res, err := q.MessageRepo.GetMessages(message)
	if err != nil {
		q.Log.Errorf("Can't get messages: %s", res)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
	out, _ := json.Marshal(res)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)

	return
}

func (q *MessageHandlers) Route(router *mux.Router) {
	s := router.PathPrefix("/messages").Subrouter()
	s.HandleFunc("/add", q.Create).Methods("POST")
	s.HandleFunc("/get", q.GetMessages).Methods("POST")
}
