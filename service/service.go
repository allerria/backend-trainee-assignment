package service

import (
	"encoding/json"
	"fmt"
	"github.com/allerria/backend-trainee-assignment/models"
	"github.com/caarlos0/env/v6"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"strings"
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

type appError struct {
	Error error
	Code  int
}

type appHandler func(http.ResponseWriter, *http.Request) *appError

func ParseConfig() (*ConfigService, error) {
	cfg := &ConfigService{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func InitService() (*Service, error) {
	db, err := models.InitDB()
	if err != nil {
		return nil, err
	}
	cfg, err := ParseConfig()
	if err != nil {
		return nil, err
	}
	s := &Service{
		Model: db,
	}
	s.Server = http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: CreateRouter(s),
	}
	return s, nil
}

func (s *Service) Serve() {
	defer s.Model.(*models.DB).Close()
	port := strings.Split(s.Server.Addr, ":")[1]
	log.Println(fmt.Sprintf("Start server on port %s", port))
	log.Fatal(s.Server.ListenAndServe())
}

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil {
		log.Printf("Error: %s %s - %s", r.Method, r.URL.EscapedPath(), e.Error)
		http.Error(w, e.Error.Error(), e.Code)
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

func (s *Service) createUserHandler(w http.ResponseWriter, r *http.Request) *appError {
	defer r.Body.Close()

	data := CreateUserRequestBody{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return &appError{err, http.StatusInternalServerError}
	}

	id, err := s.Model.CreateUser(data.Username)
	if err != nil {
		return &appError{err, http.StatusInternalServerError}
	}

	if err := json.NewEncoder(w).Encode(map[string]string{"id": id}); err != nil {
		return &appError{err, http.StatusInternalServerError}
	}
	return nil
}

func (s *Service) creatChatHandler(w http.ResponseWriter, r *http.Request) *appError {
	defer r.Body.Close()

	data := CreateChatRequestBody{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return &appError{err, http.StatusInternalServerError}
	}

	id, err := s.Model.CreateChat(data.Name, data.Users)
	if err != nil {
		return &appError{err, http.StatusInternalServerError}
	}

	if err := json.NewEncoder(w).Encode(map[string]string{"id": strconv.Itoa(int(id))}); err != nil {
		return &appError{err, http.StatusInternalServerError}
	}
	return nil
}

func (s *Service) createMessageHandler(w http.ResponseWriter, r *http.Request) *appError {
	defer r.Body.Close()

	data := CreateMessageRequestBody{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return &appError{err, http.StatusInternalServerError}
	}

	var chatID int
	chatID, err := strconv.Atoi(data.Chat)
	if err != nil {
		return &appError{err, http.StatusInternalServerError}
	}

	id, err := s.Model.CreateMessage(uint64(chatID), data.Author, data.Text)
	if err != nil {
		return &appError{err, http.StatusInternalServerError}
	}

	if err := json.NewEncoder(w).Encode(map[string]string{"id": strconv.Itoa(int(id))}); err != nil {
		return &appError{err, http.StatusInternalServerError}
	}
	return nil
}

func (s *Service) getUserChatsHandler(w http.ResponseWriter, r *http.Request) *appError {
	defer r.Body.Close()

	data := GetUserChatsRequestBody{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return &appError{err, http.StatusInternalServerError}
	}

	chats, err := s.Model.GetUserChats(data.ID)
	if err != nil {
		return &appError{err, http.StatusInternalServerError}
	}

	if err := json.NewEncoder(w).Encode(chats); err != nil {
		return &appError{err, http.StatusInternalServerError}
	}
	return nil
}

func (s *Service) getChatMessagesHandler(w http.ResponseWriter, r *http.Request) *appError {
	defer r.Body.Close()

	data := GetChatMessagesRequestBody{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return &appError{err, http.StatusInternalServerError}
	}

	var chatID int
	chatID, err := strconv.Atoi(data.Chat)
	if err != nil {
		return &appError{err, http.StatusInternalServerError}
	}

	messages, err := s.Model.GetChatMessages(uint64(chatID))
	if err != nil {
		return &appError{err, http.StatusInternalServerError}
	}

	if err := json.NewEncoder(w).Encode(messages); err != nil {
		return &appError{err, http.StatusInternalServerError}
	}
	return nil
}
