package main

import (
	"net/http"

	"github.com/Reywaltz/avito_backend/cmd/avito_api/handlers"
	"github.com/Reywaltz/avito_backend/configs"
	"github.com/Reywaltz/avito_backend/internal/repository/chat_repo"
	"github.com/Reywaltz/avito_backend/internal/repository/message_repo"
	"github.com/Reywaltz/avito_backend/internal/repository/user_repo"
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

	cfg, err := configs.InitConfig()
	if err != nil {
		log.Fatalf("Can't init config object")
	}

	db, err := postgres.NewDB(cfg)
	if err != nil {
		log.Fatalf("Can't connect to database: %s", err.Error())
	}

	user_rep := user_repo.NewUserRepository(db, log)
	log.Infof("Created UserRepo")
	chat_rep := chat_repo.NewChatRepository(db, log)
	log.Infof("Created ChatRepo")
	message_rep := message_repo.NewMessageRepository(db, log)
	log.Infof("Created ChatRepo")

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

	http.Handle("/", router)

	log.Infof("Server is up")

	http.ListenAndServe(":9000", router)
}
