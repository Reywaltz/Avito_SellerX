package main

import (
	"net/http"

	"github.com/Reywaltz/avito_backend/cmd/avito_api/handlers"
	chatRepo "github.com/Reywaltz/avito_backend/internal/repository/chat"
	messageRepo "github.com/Reywaltz/avito_backend/internal/repository/message"
	userRepo "github.com/Reywaltz/avito_backend/internal/repository/user"
	log "github.com/Reywaltz/avito_backend/pkg/log"
	"github.com/Reywaltz/avito_backend/pkg/postgres"
	"github.com/gorilla/mux"
)

func main() {
	log, err := log.NewLogger()
	if err != nil {
		log.Fatalf("Can't create logger: %s", err.Error())
	}
	log.Infof("Inited Logger")

	db, err := postgres.NewDB()
	if err != nil {
		log.Fatalf("Can't connect to database: %s", err.Error())
	}

	user_rep := userRepo.NewUserRepository(db)
	log.Infof("Created UserRepo")
	chat_rep := chatRepo.NewChatRepository(db)
	log.Infof("Created ChatRepo")
	message_rep := messageRepo.NewMessageRepository(db)
	log.Infof("Created MessageRepo")

	userHandlers := handlers.NewUserHandlers(log, user_rep)
	log.Infof("Created User handlers")
	chatHandlers := handlers.NewChatHandlers(log, chat_rep, user_rep)
	log.Infof("Created Chat handlers")
	messageHandlers := handlers.NewMessageHandlers(log, message_rep)
	log.Infof("Created Message handlers")

	router := mux.NewRouter()
	router.StrictSlash(true)

	userHandlers.Route(router)
	chatHandlers.Route(router)
	messageHandlers.Route(router)
	log.Infof("Inited router mux")

	http.Handle("/", router)

	log.Infof("Server is up")

	http.ListenAndServe(":9000", router)
}
