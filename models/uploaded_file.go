package models

import (
	"time"
)

type UploadedFile struct {
    ID         uint `gorm:"primaryKey"`
    Filename   string
    User       string
    SizeBytes  float64
    SizeKB     float64
    SizeMB     float64
    FilePath   string
    Checksum   string `gorm:"uniqueIndex"`
    UploadedAt time.Time
}