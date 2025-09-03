package jobs

import (
	"fmt"
	"log"
	"time"

	"github.com/MyEcoSys/document-uploader/database"
	"github.com/MyEcoSys/document-uploader/models"
	"github.com/MyEcoSys/document-uploader/utils"
)

func ProcessUploadTask(task models.Task) {
    // Compute SHA256 checksum
    checksum, err := utils.ComputeSHA256(task.FilePath)
    if err != nil {
        log.Printf("Failed to compute checksum: %v", err)
        return
    }
	log.Printf("File processed: %s, Checksum: %s", task.Filename, checksum)

    // Store metadata in DB
    record := models.UploadedFile{
        Filename:   task.Filename,
        User:       task.UploadedBy,
        Checksum:   checksum,
        UploadedAt: time.Now(),
    }

    err = database.SaveUploadedFile(record)
    if err != nil {
        log.Printf("DB insert failed: %v", err)
    } else {
        fmt.Println("Upload processed:", record.Filename)
    }
}
