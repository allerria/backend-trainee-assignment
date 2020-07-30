package main

import (
	"fmt"
	"github.com/allerria/backend-trainee-assignment/models"
	"github.com/allerria/backend-trainee-assignment/service"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hi!")
}

func serve() {
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)
	log.Fatal(http.ListenAndServe(":8080", r))
}

// Для дебага на windows
func setEnvVariables() {
	os.Setenv("POSTGRES_USER", "allerria")
	os.Setenv("POSTGRES_PASS", "root")
	os.Setenv("POSTGRES_DATABASE", "messenger")
}

func main() {
	setEnvVariables()
	cfgDB, err := models.ParseConfigDB()
	if err != nil {
		log.Fatal(err)
	}
	db, err := models.InitDB(cfgDB)
	if err != nil {
		log.Fatal(err)
	}
	s := &service.Service{
		DB: db,
	}
	s.Server = http.Server{
		Addr:    ":9000",
		Handler: service.CreateRouter(s),
	}
	log.Fatal(s.Server.ListenAndServe())
}
