package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Reywaltz/avito_backend/internal/models/messages"
	log "github.com/Reywaltz/avito_backend/pkg/log"
	"github.com/gorilla/mux"
)

type MessageRepository interface {
	Create(message messages.Message) (int, error)
	GetMessages(message messages.Message) ([]messages.Message, error)
	GetChatMessages(message messages.Message) ([]messages.Message, error)
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
	var message messages.Message

	if err := message.PostBind(r); err != nil {
		q.Log.Errorf("Can't insert message: %s", err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	messageID, err := q.MessageRepo.Create(message)
	if err != nil {
		q.Log.Errorf("Can't create new message: %s", err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	w.WriteHeader(http.StatusCreated)
	q.Log.Infof("New id: %d", messageID)
}

func (q *MessageHandlers) GetMessages(w http.ResponseWriter, r *http.Request) {
	var message messages.Message

	if err := message.GetBind(r); err != nil {
		q.Log.Errorf("Error bind: %s", err)
		w.WriteHeader(http.StatusBadRequest)

		return
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

func (q *MessageHandlers) GetChatMessages(w http.ResponseWriter, r *http.Request) {
	var message messages.Message

	if err := message.GetBind(r); err != nil {
		q.Log.Errorf("Can't bind Json: %s", err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	res, err := q.MessageRepo.GetChatMessages(message)
	if err != nil {
		q.Log.Errorf("Can't get messages: %s", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	out, err := json.Marshal(res)
	if err != nil {
		q.Log.Errorf("Can't marshall res: %s", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

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
