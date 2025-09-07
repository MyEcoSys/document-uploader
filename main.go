package main

import (
	"github.com/MyEcoSys/document-uploader/config"
	"github.com/MyEcoSys/document-uploader/database"
	"github.com/MyEcoSys/document-uploader/server"
)


func main() {
	// Load configuration
	//client, err := gdrive.NewGDriveClient("google_client_secret.json")
	//if err != nil {
	//	log.Fatalf("Failed to initialize Google Drive client: %v", err)
	//}
	//_ = client // Use the client as needed

	//client.UploadFile("test.txt", "text/plain", "root")

	config.LoadConfig()
	
	database.InitDB()
	
	server := server.NewServer()
	server.Start()
}
