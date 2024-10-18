package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func downloadZipHandler(w http.ResponseWriter, r *http.Request) {
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

	zipPath := filepath.Join(uploadDir, folder)
	zipName := zipPath + ".zip"
	_, err := os.Stat(zipName)
	if err == nil {
		response.StatusCode = http.StatusOK
		response.Status = "OK"
		response.Message = "Zipping downloaded"
		response.ResponseTime = time.Now()
		jsonResponse.Encode(response)
		http.ServeFile(w, r, zipName)
		return
	}
	
	
	status, ok := zipStatuses.Load(folder)
	if !ok {
		http.Error(w, "No zipping found for the specified folder", http.StatusNotFound)
		response.StatusCode = http.StatusNotFound
		response.Status = "Not Found"
		response.Message = "No zipping found for the specified folder"
		response.ResponseTime = time.Now()
		jsonResponse.Encode(response)
		return
	} else if status.(*zipStatus).Status != "completed" {
		http.Error(w, "Zipping not completed or in progress", http.StatusConflict)
		response.StatusCode = http.StatusConflict
		response.Status = "Conflict"
		response.Message = "Zipping not completed or in progress"
		response.ResponseTime = time.Now()
		jsonResponse.Encode(response)
		return
	} else {
		response.StatusCode = http.StatusOK
		response.Status = "OK"
		response.Message = "Zipping downloaded"
		response.ResponseTime = time.Now()
		zipFilePath := status.(*zipStatus).FilePath
		http.ServeFile(w, r, zipFilePath)
		log.Println(response)
		jsonResponse.Encode(response)
	}
}
