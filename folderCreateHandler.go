package main

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func folderCreateHandler(w http.ResponseWriter, r *http.Request) {
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
	folderInfo, err := os.Stat(folderPath)
	if err != nil {
		os.MkdirAll(folderPath, os.ModePerm)
		response.StatusCode = http.StatusCreated
		response.Status = "Created"
		response.Message = "Folder created successfully"
		response.ResponseTime = time.Now()
		jsonResponse.Encode(response)
	} else {
		response.StatusCode = http.StatusCreated
		response.Status = "Created"
		response.Message = "Folder created successfully"
		response.ResponseTime = time.Now()
		jsonResponse.Encode(response)
	}
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
}
