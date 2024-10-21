package main

import (
	"fmt"
	"io/fs"
	"net/http"
	"sync"
	"time"
)

const (
	uploadDir          = "./uploads"
	fileMetadataFile   = "./fileMetadata.json"
	folderMetadataFile = "./folderMetadata.json"
	workerCount        = 100
)

var (
	fileLocks         sync.Map
	metaDataMap       sync.Map
	metadataLock      sync.Mutex
	folderMetadataMap sync.Map
	taskQueue         = make(chan string, 150)
	wg                sync.WaitGroup
	zipStatuses       sync.Map
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
	FolderMode   fs.FileMode `json:"fileMode"`
	IsDirectory  bool        `json:"is_directory"`
}

type zipStatus struct {
	Status    string
	StartTime time.Time
	EndTime   time.Time
	FilePath  string
	ErrorMsg  string
}

func main() {
	loadFileMetadata()
	loadFolderMetadata()
	workerPool()
	http.HandleFunc("/upload", uploadFileHandler)
	http.HandleFunc("/download", downloadFileHandler)
	http.HandleFunc("/fileinfo", fileInfoHandler)
	http.HandleFunc("/delete", deleteFileHandler)
	http.HandleFunc("/folderinfo", folderInfoHandler)
	http.HandleFunc("/createfolder", folderCreateHandler)
	http.HandleFunc("/deletefolder", deleteFolderHandler)
	http.HandleFunc("/zip", zipFolderHandler)
	http.HandleFunc("/zipdownload", downloadZipHandler)
	http.HandleFunc("/zipstatus", zipStatusHandler)
	fmt.Println("Server is listening in http://localhost:8080")
	http.ListenAndServe(":8080", nil)
	Shutdown()
}