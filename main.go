package main

import (
	"github.com/allerria/backend-trainee-assignment/service"
	"log"
	"os"
)

// Для дебага на windows
func setEnvVariables() {
	os.Setenv("POSTGRES_USER", "allerria")
	os.Setenv("POSTGRES_PASS", "root")
	os.Setenv("POSTGRES_DATABASE", "messenger")
	os.Setenv("SERVICE_PORT", "9000")
}

func main() {
	setEnvVariables()
	service, err := service.InitService()
	if err != nil {
		log.Fatal(err)
	}
	service.Serve()
}
