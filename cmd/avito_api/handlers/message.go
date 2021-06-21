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

func (q *MessageHandlers) Create(writer http.ResponseWriter, request *http.Request) {
	type incomeJson map[string]string

	var message messages.Message

	err := message.PostBind(request)
	if err != nil {
		q.Log.Errorf("Can't insert message: %s", err)
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	messageID, err := q.MessageRepo.Create(message)
	if err != nil {
		q.Log.Errorf("Can't create new message: %s", err)
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	writer.WriteHeader(http.StatusCreated)
	q.Log.Infof("New id: %d", messageID)
}

func (q *MessageHandlers) GetMessages(writer http.ResponseWriter, request *http.Request) {
	var message messages.Message

	err := message.GetBind(request)
	if err != nil {
		q.Log.Errorf("Error bind: %s", err)
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	res, err := q.MessageRepo.GetMessages(message)
	if err != nil {
		q.Log.Errorf("Can't get messages: %s", res)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	out, _ := json.Marshal(res)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(out)

	return
}

func (q *MessageHandlers) GetChatMessages(writer http.ResponseWriter, request *http.Request) {
	var message messages.Message

	if err := message.GetBind(request); err != nil {
		q.Log.Errorf("Can't bind Json: %s", err)

		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	res, err := q.MessageRepo.GetChatMessages(message)
	if err != nil {
		q.Log.Errorf("Can't get messages: %s", err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	out, err := json.Marshal(res)
	if err != nil {
		q.Log.Errorf("Can't marshall res: %s", err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(out)

	return
}

func (q *MessageHandlers) Route(router *mux.Router) {
	s := router.PathPrefix("/messages").Subrouter()
	s.HandleFunc("/add", q.Create).Methods("POST")
	s.HandleFunc("/get", q.GetMessages).Methods("POST")
}
