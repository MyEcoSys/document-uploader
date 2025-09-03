package models

import "time"

type UploadedFile struct {
    ID         uint `gorm:"primaryKey"`
    Filename   string
    User       string
    Checksum   string
    UploadedAt time.Time
}
