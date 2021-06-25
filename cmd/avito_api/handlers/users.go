package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Reywaltz/avito_backend/internal/additions"
	"github.com/Reywaltz/avito_backend/internal/message"
	"github.com/Reywaltz/avito_backend/internal/models/users"
	log "github.com/Reywaltz/avito_backend/pkg/log"
	"github.com/Reywaltz/avito_backend/pkg/postgres"
	"github.com/gorilla/mux"
)

type UserRepository interface {
	Create(user users.User) (int, error)
	GetAll() ([]users.User, error)
	GetOne(ID int) (users.User, error)
}

type UserHandlers struct {
	Log      log.Logger
	UserRepo UserRepository
}

func NewUserHandlers(logger log.Logger, userRepo UserRepository) *UserHandlers {
	return &UserHandlers{
		Log:      logger,
		UserRepo: userRepo,
	}
}

func (q *UserHandlers) GetAll(w http.ResponseWriter, r *http.Request) {
	out, _ := q.UserRepo.GetAll()

	res, _ := json.Marshal(out)
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func (q *UserHandlers) Create(w http.ResponseWriter, r *http.Request) {
	var user users.User

	if err := user.Bind(r); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		q.Log.Errorf("Got emplty json")

		return
	}

	user.CreatedAt = time.Now()

	if additions.ValidateUser(user) {
		id, err := q.UserRepo.Create(user)
		if err != nil {
			if errors.Is(err, postgres.DuplicateError) {
				w.WriteHeader(http.StatusBadRequest)
				q.Log.Errorf("Can't create entity: %s", err)

				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			q.Log.Errorf("Can't create entity: %s", err)

			return
		}

		message.MakeResponse(w, id, http.StatusCreated)
		q.Log.Infof("Created user with id: {%d}", id)

		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
		q.Log.Errorf("Wrong JSON input")

		return
	}
}

func (q *UserHandlers) Route(router *mux.Router) {
	s := router.PathPrefix("/users").Subrouter()
	s.HandleFunc("/", q.GetAll).Methods("GET")
	s.HandleFunc("/add", q.Create).Methods("POST")
}
