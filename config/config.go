package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort       string
	UploadDir        string
	MaxUploadSizeMB  int
	AllowedFileTypes []string
}

var AppConfig Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Using default config.")
	}
	// For simplicity, hardcoding values. In a real application, load from environment variables or config files.
	AppConfig = Config{
		ServerPort:       getEnv("SERVER_PORT", "8080"),
		UploadDir:        getEnv("UPLOAD_DIR", "./uploads"),
		MaxUploadSizeMB:  getEnvAsInt("MAX_UPLOAD_SIZE_MB", 50),
		AllowedFileTypes: strings.Split(getEnv("ALLOWED_FILE_TYPES", "application/pdf,text/csv,application/vnd.openxmlformats-officedocument.wordprocessingml.document"), ","),
	}
	log.Printf("Config loaded: %+v\n", AppConfig)
}

func getEnv(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(name string, defaultVal int) int {
	if valueStr := os.Getenv(name); valueStr != "" {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultVal
}
