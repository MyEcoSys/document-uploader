package database

import (
    "gorm.io/gorm"
    "github.com/MyEcoSys/document-uploader/models"
)

var DB *gorm.DB

func InitDB(db *gorm.DB) {
    DB = db
    DB.AutoMigrate(&models.UploadedFile{})
}

func SaveUploadedFile(f models.UploadedFile) error {
    return DB.Create(&f).Error
}