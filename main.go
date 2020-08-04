package main

import (
	"github.com/allerria/backend-trainee-assignment/service"
	"log"
)

func main() {
	s, err := service.InitService()
	if err != nil {
		log.Fatal(err)
	}
	s.Serve()
}
