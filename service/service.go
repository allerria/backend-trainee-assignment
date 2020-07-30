package service

import (
	"encoding/json"
	"github.com/allerria/backend-trainee-assignment/models"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type Service struct {
	DB     *models.DB
	Server http.Server
}

type CreateUserRequestBody struct {
	Username string `json:"username"`
}

func CreateRouter(s *Service) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/users", s.createUserHandler).Methods(http.MethodPost)
	return r
}

func (s *Service) createUserHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data := CreateUserRequestBody{}
	if err := json.Unmarshal(body, &data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	id, err := s.DB.CreateUser(data.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp := struct {
		ID string `json:"id"`
	}{id}
	output, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}
