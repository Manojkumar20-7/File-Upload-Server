package config

import (
	"io/fs"
	"os"
	"time"
)


type Response struct {
	StatusCode   int       `json:"status_code"`
	Status       string    `json:"status"`
	Message      string    `json:"message"`
	ResponseTime time.Time `json:"response_time"`
}

type FileMetadata struct {
	FileName     string      `json:"file_name"`
	FilePath     string      `json:"file_path"`
	FolderPath   string      `json:"folder_path"`
	FileSize     int64       `json:"file_size"`
	ModifiedTime string      `json:"modified_time"`
	CreatedTime  string      `json:"created_time"`
	FileMode     fs.FileMode `json:"file_mode"`
	IsDirectory  bool        `json:"is_directory"`
}

type FolderMetadata struct {
	FolderName   string      `json:"folder_name"`
	FolderPath   string      `json:"folder_path"`
	FolderSize   int64       `json:"folder_size"`
	FilesCount   int         `json:"files_count"`
	ModifiedTime string      `json:"modified_time"`
	CreatedTime  string      `json:"created_time"`
	FolderMode   os.FileMode `json:"fileMode"`
	IsDirectory  bool        `json:"is_directory"`
}

type ZipStatus struct {
	Status    string
	StartTime time.Time
	EndTime   time.Time
	FilePath  string
	ErrorMsg  string
}