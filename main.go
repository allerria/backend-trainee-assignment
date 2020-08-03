package main

import (
	"github.com/allerria/backend-trainee-assignment/service"
	"log"
)

func main() {
	service, err := service.InitService()
	if err != nil {
		log.Fatal(err)
	}
	service.Serve()
}
