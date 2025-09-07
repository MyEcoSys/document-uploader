package server

import (
	"fmt"
	"io"
	"log"
	"math"
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
	s.SetupRoutes()
	queue.StartWorkerPool(3)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      s.router,
		ReadTimeout:  15 * time.Minute,  // max time to read entire request, including body
		WriteTimeout: 15 * time.Minute,  // max time to write response
		IdleTimeout:  60 * time.Second,  // max time to keep idle connections alive
	}
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
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

	fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filename)
	filePath := filepath.Join(config.AppConfig.UploadDir, fileName)

	// Ensure the uploads directory exists
	if err := ensureDir(filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	// Open a new file for writing
	outFile, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file"})
		return
	}
	defer outFile.Close()

	errWrite := make(chan error, 1)

	go func() {
		written, err := io.Copy(outFile, c.Request.Body)
		if err != nil {
			errWrite <- err
		}
		if written == 0 {
			errWrite <- fmt.Errorf("empty file")
		}
		errWrite <- nil
	}()

	// Wait for either the write to finish or a timeout
	select {
	case err := <-errWrite:
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}
	case <-time.After(5 * time.Minute):
		c.JSON(http.StatusRequestTimeout, gin.H{"error": "File upload timed out"})
		return
	}

	stat, _ := outFile.Stat()
	sizeBytes := stat.Size()
	sizeKB := math.Round(float64(sizeBytes) / 1024)
	sizeMB := math.Round(float64(sizeBytes) / (1024 * 1024))

	queue.Enqueue(models.Task{
		FilePath:  filePath,
		Filename:  filename,
		UploadedBy: "anonymous", // Replace with actual user info
		SizeBytes: float64(sizeBytes),
		SizeKB:    sizeKB,
		SizeMB:    sizeMB,
	})

	c.JSON(http.StatusAccepted, gin.H{
		"message":  "File uploaded successfully",
		"filename": fileName,
		"path":     filePath,
		"fileSize": gin.H{
			"bytes": sizeBytes,
			"kb":    sizeKB,
			"mb":    sizeMB,
		},
	})
}


func ensureDir(filePath string) error {
	dir := filepath.Dir(filePath)
	return os.MkdirAll(dir, os.ModePerm)
}