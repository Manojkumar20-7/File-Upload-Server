package main

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
)

func folderCreateHandler(w http.ResponseWriter, r *http.Request) {
	logField := log.Fields{
		"method": "folderCreateHandler",
	}
	logger.Log(log.InfoLevel, logField, "Folder create handler begin")
	response := Response{}
	jsonResponse := json.NewEncoder(w)
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
		logger.Log(log.TraceLevel, logField, "Creating folder")
		os.MkdirAll(folderPath, os.ModePerm)
		response.StatusCode = http.StatusCreated
		response.Status = "Created"
		response.Message = "Folder created successfully"
		response.ResponseTime = time.Now()
		jsonResponse.Encode(response)
	} else {
		logger.Log(log.ErrorLevel, logField, "Error in creating folder")
		response.StatusCode = http.StatusCreated
		response.Status = "Created"
		response.Message = "Folder created successfully"
		response.ResponseTime = time.Now()
		jsonResponse.Encode(response)
	}
	logger.Log(log.DebugLevel, logField, "Updating folder metadata")
	folderInfo, _ := os.Stat(folderPath)
	folderMetadata := FolderMetadata{}
	folderMetadata.FolderName = folderInfo.Name()
	folderMetadata.FolderPath = folderPath
	folderMetadata.FolderSize = folderInfo.Size()
	folderMetadata.FilesCount = 0
	folderMetadata.ModifiedTime = folderInfo.ModTime().Format(http.TimeFormat)
	folderMetadata.CreatedTime = time.Now().Format(http.TimeFormat)
	folderMetadata.FolderMode = folderInfo.Mode()
	folderMetadata.IsDirectory = folderInfo.IsDir()
	folderMetadataMap.Store(folderPath, folderMetadata)
	saveFolderMetadata()
	logger.Log(log.InfoLevel, logField, "Completing folder create and exits")
}
