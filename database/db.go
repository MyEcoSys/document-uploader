package database

import (
	"fmt"
	"log"
	"sync"

	"github.com/MyEcoSys/document-uploader/config"
	"github.com/MyEcoSys/document-uploader/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB
var mu sync.Mutex

func InitDB() {
    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
        config.AppConfig.DbHost,
        config.AppConfig.DbUser,
        config.AppConfig.DbPassword,
        config.AppConfig.DbName,
        config.AppConfig.DbPort,
    )
    db, err := gorm.Open(postgres.New(postgres.Config{
        DSN:                  dsn,
        PreferSimpleProtocol: true,
        DriverName:           "pgx", // Use pgx driver
    }), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })

    if err != nil {
        log.Fatalf("failed to connect to database: %v", err)
    }

    DB = db
    if err := DB.AutoMigrate(&models.UploadedFile{}); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}
}

func SaveUploadedFile(f models.UploadedFile) error {
    mu.Lock()
    defer mu.Unlock()

    // Check for duplicate checksum

    var existing models.UploadedFile
    if err := DB.Where("checksum = ?", f.Checksum).First(&existing).Error; err == nil {
        return fmt.Errorf("duplicate file: checksum already exists")
    }

    return DB.Create(&f).Error
}

func GetUploadedFileByID(id uint) (*models.UploadedFile, error) {
    var file models.UploadedFile
    if err := DB.First(&file, id).Error; err != nil {
        return nil, err
    }
    return &file, nil
}