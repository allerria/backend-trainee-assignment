package service

import (
	"encoding/json"
	"github.com/allerria/backend-trainee-assignment/models"
	"github.com/caarlos0/env/v6"
	"github.com/gorilla/mux"
	"io/ioutil"
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

func CreateRouter(s *Service) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/users/add", s.createUserHandler).Methods(http.MethodPost)
	r.HandleFunc("/chats/add", s.creatChatHandler).Methods(http.MethodPost)
	r.HandleFunc("/chats/get", s.getUserChatsHandler).Methods(http.MethodPost)
	r.HandleFunc("/messages/add", s.createMessageHandler).Methods(http.MethodPost)
	r.HandleFunc("/messages/get", s.getChatMessagesHandler).Methods(http.MethodPost)
	return r
}

func ParseConfig() (*ConfigService, error) {
	cfg := &ConfigService{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (s *Service) createUserHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		msg, _ := json.Marshal(ErrorResponse{"Error! Can't read request body."})
		w.Write(msg)
		return
	}

	data := CreateUserRequestBody{}
	if err := json.Unmarshal(body, &data); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		msg, _ := json.Marshal(ErrorResponse{"Error! Can't unmarshal request."})
		w.Write(msg)
		return
	}

	id, err := s.Model.CreateUser(data.Username)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		msg, _ := json.Marshal(ErrorResponse{"Error! Can't create user."})
		w.Write(msg)
		return
	}

	msg, err := json.Marshal(map[string]string{"id": id})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		msg, _ := json.Marshal(ErrorResponse{"User have been created, but can't marshal id."})
		w.Write(msg)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(msg)
}

func (s *Service) creatChatHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		msg, _ := json.Marshal(ErrorResponse{"Error! Can't read request body."})
		w.Write(msg)
		return
	}

	data := CreateChatRequestBody{}
	if err := json.Unmarshal(body, &data); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		msg, _ := json.Marshal(ErrorResponse{"Error! Can't unmarshal request."})
		w.Write(msg)
		return
	}

	id, err := s.Model.CreateChat(data.Name, data.Users)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		msg, _ := json.Marshal(ErrorResponse{"Error! Can't create chat."})
		w.Write(msg)
		return
	}

	msg, err := json.Marshal(map[string]uint64{"id": id})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		msg, _ := json.Marshal(ErrorResponse{"Chat created, but can't marshal id."})
		w.Write(msg)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(msg)
}

func (s *Service) createMessageHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		msg, _ := json.Marshal(ErrorResponse{"Error! Can't read request body."})
		w.Write(msg)
		return
	}

	data := CreateMessageRequestBody{}
	if err := json.Unmarshal(body, &data); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		msg, _ := json.Marshal(ErrorResponse{"Error! Can't unmarshal request."})
		w.Write(msg)
		return
	}
	var chatID int
	chatID, err = strconv.Atoi(data.Chat)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		msg, _ := json.Marshal(ErrorResponse{"Error! ChatID is incorrect."})
		w.Write(msg)
		return
	}

	id, err := s.Model.CreateMessage(uint64(chatID), data.Author, data.Text)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		msg, _ := json.Marshal(ErrorResponse{"Error! Can't create message."})
		w.Write(msg)
		return
	}

	msg, err := json.Marshal(map[string]uint64{"id": id})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		msg, _ := json.Marshal(ErrorResponse{"Message created, but can't unmarshal id."})
		w.Write(msg)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(msg)
}

func (s *Service) getUserChatsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		msg, _ := json.Marshal(ErrorResponse{"Error! Can't read request body."})
		w.Write(msg)
		return
	}

	data := GetUserChatsRequestBody{}
	if err := json.Unmarshal(body, &data); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		msg, _ := json.Marshal(ErrorResponse{"Error! Can't unmarshal request."})
		w.Write(msg)
		return
	}

	chats, err := s.Model.GetUserChats(data.ID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		msg, _ := json.Marshal(ErrorResponse{"Error! Can't get user chats."})
		w.Write(msg)
		return
	}

	msg, err := json.Marshal(chats)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		msg, _ := json.Marshal(ErrorResponse{"Error! Can't marshal chats."})
		w.Write(msg)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(msg)
}

func (s *Service) getChatMessagesHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		msg, _ := json.Marshal(ErrorResponse{"Error! Can't read request body."})
		w.Write(msg)
		return
	}

	data := GetChatMessagesRequestBody{}
	if err := json.Unmarshal(body, &data); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		msg, _ := json.Marshal(ErrorResponse{"Error! Can't unmarshal request."})
		w.Write(msg)
		return
	}
	var chatID int
	chatID, err = strconv.Atoi(data.Chat)
	if err := json.Unmarshal(body, &data); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		msg, _ := json.Marshal(ErrorResponse{"Error! ChatID is incorrect."})
		w.Write(msg)
		return
	}

	chats, err := s.Model.GetChatMessages(uint64(chatID))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		msg, _ := json.Marshal(ErrorResponse{"Error! Can't get messages."})
		w.Write(msg)
		return
	}

	msg, err := json.Marshal(chats)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		msg, _ := json.Marshal(ErrorResponse{"Error! Can't marshal messages."})
		w.Write(msg)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(msg)
}
