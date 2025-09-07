package models

type Task struct {
	FilePath   string
	Filename   string
	UploadedBy string
	SizeBytes  float64
	SizeKB     float64
	SizeMB     float64
}
