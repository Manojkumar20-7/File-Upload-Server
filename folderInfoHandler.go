package main

import (
	"encoding/json"
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

func folderInfoHandler(w http.ResponseWriter, r *http.Request) {
	logField := log.Fields{
		"method": "folderInfoHandler",
	}
	logger.Log(log.InfoLevel, logField, "Folder info handler begins")
	response := Response{}
	jsonResponse := json.NewEncoder(w)
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
	_, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		http.Error(w, "Folder not found", http.StatusNotFound)
		response.StatusCode = http.StatusNotFound
		response.Status = "Not Found"
		response.Message = "Folder not found"
		response.ResponseTime = time.Now()
		jsonResponse.Encode(response)
		return
	}
	if err != nil {
		http.Error(w, "Error in reading folder info", http.StatusNotFound)
		response.StatusCode = http.StatusNotFound
		response.Status = "Not Found"
		response.Message = "Folder not found"
		response.ResponseTime = time.Now()
		jsonResponse.Encode(response)
		return
	}
	var result any
	folderMetadataMap.Range(func(key, value any) bool {
		if key == folderPath {
			result = value
		}
		return true
	})
	response.StatusCode = http.StatusOK
	response.Status = "OK"
	response.Message = "Folder info retrieved successfully"
	response.ResponseTime = time.Now()
	jsonResponse.Encode(response)
	json.NewEncoder(w).Encode(result)
	logger.Log(log.InfoLevel, logField, "Folder info handler exits")
}
