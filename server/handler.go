package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/MyEcoSys/document-uploader/config"
	"github.com/MyEcoSys/document-uploader/models"
	"github.com/MyEcoSys/document-uploader/queue"
	"github.com/gin-gonic/gin"
)

type Server struct{
	router *gin.Engine
}

func NewServer() *Server {
	return &Server{
		router: gin.Default(),
	}
}

func (s *Server) Start() error{
	log.Println("Server started successfully on port 8080")
	queue.StartWorkerPool(3)
	log.Fatal(http.ListenAndServe(":8080", s.router))
	return nil
}

func (s *Server) SetupRoutes() {
	s.router.POST("/login", loginHandler)
	s.router.POST("/upload", s.uploadHandler)
}

func loginHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Login endpoint"})
}


func (s *Server) uploadHandler(c *gin.Context) {
	maxSize := int64(config.AppConfig.MaxUploadSizeMB * 1024 * 1024)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)

	filename := c.GetHeader("X-Filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}
	log.Printf("Received file: %s", filename)

	// Read all raw data
    data, err := io.ReadAll(c.Request.Body)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
        return
    }

	if len(data) == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Empty file"})
        return
    }

	fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filename)
	filePath := filepath.Join(config.AppConfig.UploadDir, fileName)
	
	// Ensure the uploads directory exists
	if err := ensureDir(filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	queue.Enqueue(models.Task{
        FilePath:  filePath,
        Filename:  filename,
        UploadedBy: "anonymous", // Replace with actual user info
    })
	c.JSON(http.StatusAccepted, gin.H{"message": "File uploaded successfully", "filename": fileName})
}

func ensureDir(filePath string) error {
	dir := filepath.Dir(filePath)
	return os.MkdirAll(dir, os.ModePerm)
}