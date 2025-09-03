package main

import (
	"log"

	"github.com/MyEcoSys/document-uploader/config"
	"github.com/MyEcoSys/document-uploader/server"
)


func main() {
	config.LoadConfig()
	server := server.NewServer()

	server.SetupRoutes()
	server.Start()
	log.Println("Server started")
}