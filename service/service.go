package service

import (
	"encoding/json"
	"github.com/allerria/backend-trainee-assignment/models"
	"github.com/caarlos0/env/v6"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type ConfigService struct {
	Port string `env:"SERVICE_PORT" envDefault:"9000"`
}

type Service struct {
	Model  models.Model
	Server http.Server
}

type CreateUserRequestBody struct {
	Username string `json:"username"`
}

type CreateChatRequestBody struct {
	Name  string   `json:"name"`
	Users []string `json:"users"`
}

type CreateMessageRequestBody struct {
	Chat   string `json:"chat"`
	Author string `json:"author"`
	Text   string `json:"text"`
}

type GetUserChatsRequestBody struct {
	ID string `json:"user"`
}

type GetChatMessagesRequestBody struct {
	Chat string `json:"chat"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type appError struct {
	Error   error
	Message string
	Code    int
}

type appHandler func(http.ResponseWriter, *http.Request) *appError

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil {
		log.Println(e.Error)
		http.Error(w, e.Message, e.Code)
	}
}

func CreateRouter(s *Service) *mux.Router {
	r := mux.NewRouter()
	ur := r.PathPrefix("/users").Subrouter()
	ur.HandleFunc("/add", appHandler(s.createUserHandler).ServeHTTP).Methods(http.MethodPost)
	cr := r.PathPrefix("/chats").Subrouter()
	cr.HandleFunc("/add", appHandler(s.creatChatHandler).ServeHTTP).Methods(http.MethodPost)
	cr.HandleFunc("/get", appHandler(s.getUserChatsHandler).ServeHTTP).Methods(http.MethodPost)
	mr := r.PathPrefix("/messages").Subrouter()
	mr.HandleFunc("/add", appHandler(s.createMessageHandler).ServeHTTP).Methods(http.MethodPost)
	mr.HandleFunc("/get", appHandler(s.getChatMessagesHandler).ServeHTTP).Methods(http.MethodPost)
	return r
}

func ParseConfig() (*ConfigService, error) {
	cfg := &ConfigService{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (s *Service) createUserHandler(w http.ResponseWriter, r *http.Request) *appError {
	defer r.Body.Close()

	data := CreateUserRequestBody{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return &appError{err, err.Error(), http.StatusInternalServerError}
	}

	id, err := s.Model.CreateUser(data.Username)
	if err != nil {
		return &appError{err, err.Error(), http.StatusInternalServerError}
	}

	if err := json.NewEncoder(w).Encode(map[string]string{"id": id}); err != nil {
		return &appError{err, err.Error(), http.StatusInternalServerError}
	}
	return nil
}

func (s *Service) creatChatHandler(w http.ResponseWriter, r *http.Request) *appError {
	defer r.Body.Close()

	data := CreateChatRequestBody{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return &appError{err, err.Error(), http.StatusInternalServerError}
	}

	id, err := s.Model.CreateChat(data.Name, data.Users)
	if err != nil {
		return &appError{err, err.Error(), http.StatusInternalServerError}
	}

	if err := json.NewEncoder(w).Encode(map[string]uint64{"id": id}); err != nil {
		return &appError{err, err.Error(), http.StatusInternalServerError}
	}
	return nil
}

func (s *Service) createMessageHandler(w http.ResponseWriter, r *http.Request) *appError {
	defer r.Body.Close()

	data := CreateMessageRequestBody{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return &appError{err, err.Error(), http.StatusInternalServerError}
	}

	var chatID int
	chatID, err := strconv.Atoi(data.Chat)
	if err != nil {
		return &appError{err, err.Error(), http.StatusInternalServerError}
	}

	id, err := s.Model.CreateMessage(uint64(chatID), data.Author, data.Text)
	if err != nil {
		return &appError{err, err.Error(), http.StatusInternalServerError}
	}

	if err := json.NewEncoder(w).Encode(map[string]uint64{"id": id}); err != nil {
		return &appError{err, err.Error(), http.StatusInternalServerError}
	}
	return nil
}

func (s *Service) getUserChatsHandler(w http.ResponseWriter, r *http.Request) *appError {
	defer r.Body.Close()

	data := GetUserChatsRequestBody{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return &appError{err, err.Error(), http.StatusInternalServerError}
	}

	chats, err := s.Model.GetUserChats(data.ID)
	if err != nil {
		return &appError{err, err.Error(), http.StatusInternalServerError}
	}

	if err := json.NewEncoder(w).Encode(chats); err != nil {
		return &appError{err, err.Error(), http.StatusInternalServerError}
	}
	return nil
}

func (s *Service) getChatMessagesHandler(w http.ResponseWriter, r *http.Request) *appError {
	defer r.Body.Close()

	data := GetChatMessagesRequestBody{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return &appError{err, err.Error(), http.StatusInternalServerError}
	}

	var chatID int
	chatID, err := strconv.Atoi(data.Chat)
	if err != nil {
		return &appError{err, err.Error(), http.StatusInternalServerError}
	}

	messages, err := s.Model.GetChatMessages(uint64(chatID))
	if err != nil {
		return &appError{err, err.Error(), http.StatusInternalServerError}
	}

	if err := json.NewEncoder(w).Encode(messages); err != nil {
		return &appError{err, err.Error(), http.StatusInternalServerError}
	}
	return nil
}
