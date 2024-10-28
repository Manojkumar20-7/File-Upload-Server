package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
)

type FolderDetails struct {
	FileCount    int       `json:"file_count"`
	FolderSize   int       `json:"folder_size"`
	CreatedTime  time.Time `json:"created_time"`
	ModifiedTime time.Time `json:"modified_time"`
}

func printStat(file fs.FileInfo, w http.ResponseWriter) {
	fmt.Fprintf(w, "%-30v", file.Mode())
	fmt.Fprintf(w, "%-30v", file.Name())
	fmt.Fprintf(w, "%-30v", file.Size())
	fmt.Fprintf(w, "%-30v\n", file.ModTime().Format(http.TimeFormat))
}

func folderInfoHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{}
	jsonResponse := json.NewEncoder(w)
	logField := log.Fields{
		"method": "folderInfoHandler",
	}
	logger.Log(log.InfoLevel, logField, "Folder info handler begins")
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid HTTP request", http.StatusBadRequest)
		response.StatusCode = http.StatusBadRequest
		response.Status = "Bad Request"
		response.Message = "Invalid HTTP request"
		response.ResponseTime = time.Now()
		jsonResponse.Encode(response)
		return
	}

	folder := r.URL.Query().Get("folder")
	if folder == "" {
		http.Error(w, "Folder not specified", http.StatusBadRequest)
		response.StatusCode = http.StatusBadRequest
		response.Status = "Bad Request"
		response.Message = "Folder not specified"
		response.ResponseTime = time.Now()
		jsonResponse.Encode(response)
		return
	}
	folderPath := filepath.Join(uploadDir, folder)
	folderInfo, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		http.Error(w, "Folder not found", http.StatusNotFound)
		response.StatusCode = http.StatusNotFound
		response.Status = "Not Found"
		response.Message = "Folder not found"
		response.ResponseTime = time.Now()
		jsonResponse.Encode(response)
		return
	}
	content, err := os.ReadDir(folderPath)
	if err != nil {
		http.Error(w, "Error in reading fodler info", http.StatusNotFound)
		response.StatusCode = http.StatusNotFound
		response.Status = "Not Found"
		response.Message = "File not found"
		response.ResponseTime = time.Now()
		jsonResponse.Encode(response)
		return
	}
	var result any
	metaDataMap.Range(func(key, value any) bool {
		if key == folderPath {
			result = value
		}
		return true
	})
	fmt.Fprintf(w, "Name: %s\n", folderInfo.Name())
	fmt.Fprintf(w, "Size: %d\n", folderInfo.Size())
	fmt.Fprintf(w, "Count: %d\n", len(content))
	fmt.Fprintf(w, "Last-Modified: %v\n", folderInfo.ModTime().Format(http.TimeFormat))
	filePath := filepath.Join(folderPath, content[0].Name())
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Fprintf(w, "Created-Time: %v\n", fileInfo.ModTime().Format(http.TimeFormat))
	fmt.Fprintf(w, "%-30v", "Permission")
	fmt.Fprintf(w, "%-30v", "Name")
	fmt.Fprintf(w, "%-30v", "Size")
	fmt.Fprintf(w, "%-30v\n", "Modified-time")
	for _, f := range content {
		filePath := filepath.Join(folderPath, f.Name())
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			log.Fatal(err)
			return
		}
		printStat(fileInfo, w)
	}
	json.NewEncoder(w).Encode(result)
	logger.Log(log.InfoLevel, logField, "Folder info handler exits")
}
