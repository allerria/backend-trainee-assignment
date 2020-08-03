package main

import (
	"fmt"
	"github.com/allerria/backend-trainee-assignment/models"
	"github.com/allerria/backend-trainee-assignment/service"
	"log"
	"net/http"
	"os"
)

func serve(db *models.DB, cfg *service.ConfigService) {
	s := &service.Service{
		Model: db,
	}
	s.Server = http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: service.CreateRouter(s),
	}
	log.Println("Start server on port 9000")
	log.Fatal(s.Server.ListenAndServe())
}

// Для дебага на windows
func setEnvVariables() {
	os.Setenv("POSTGRES_USER", "allerria")
	os.Setenv("POSTGRES_PASS", "root")
	os.Setenv("POSTGRES_DATABASE", "messenger")
	os.Setenv("SERVICE_PORT", "9000")
}

func main() {
	setEnvVariables()
	cfgDB, err := models.ParseConfig()
	if err != nil {
		log.Fatal(err)
	}
	db, err := models.InitDB(cfgDB)
	if err != nil {
		log.Fatal(err)
	}
	cfgService, err := service.ParseConfig()
	if err != nil {
		log.Fatal(err)
	}
	serve(db, cfgService)
}
