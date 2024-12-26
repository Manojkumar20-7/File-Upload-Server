package main

import (
	"encoding/json"
	"fileServer/config"
	"fileServer/constants"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
)

func downloadFileHandler(w http.ResponseWriter, r *http.Request) {
	logField := log.Fields{
		"method": "downloadFileHandler",
	}
	logger.Log(log.InfoLevel, logField, "Download handler begins")
	response := config.Response{}
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
	fileName := r.URL.Query().Get("filename")

	if folder == "" {
		http.Error(w, "Folder not specified", http.StatusBadRequest)
		response.StatusCode = http.StatusBadRequest
		response.Status = "Bad Request"
		response.Message = "Folder not specified"
		response.ResponseTime = time.Now()
		jsonResponse.Encode(response)
		return
	}
	if fileName == "" {
		http.Error(w, "Filename not specified", http.StatusBadRequest)
		response.StatusCode = http.StatusBadRequest
		response.Status = "Bad Request"
		response.Message = "Filename not specified"
		response.ResponseTime = time.Now()
		jsonResponse.Encode(response)
		return
	}

	filePath := filepath.Join(constants.UploadDir, folder, fileName)
	logger.Log(log.DebugLevel, logField, "Acquiring file lock")
	fileLock := getFileLock(filePath)
	fileLock.Lock()
	defer fileLock.Unlock()
	logger.Log(log.DebugLevel, logField, "Reading file for download")
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, "File not Found", http.StatusNotFound)
		response.StatusCode = http.StatusNotFound
		response.Status = "Not Found"
		response.Message = "File not Found"
		response.ResponseTime = time.Now()
		jsonResponse.Encode(response)
		return
	}
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	w.Header().Set("Content-Type", "application/octet-stream")
	response.StatusCode = http.StatusOK
	response.Status = "OK"
	response.Message = "File Downloaded successfully"
	response.ResponseTime = time.Now()
	jsonResponse.Encode(response)
	w.Write(fileContent)
	logger.Log(log.InfoLevel, logField, "Completing file download and exits")
}
