package main

import (
	"encoding/json"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func zipStatusHandler(w http.ResponseWriter, r *http.Request) {
	logField := log.Fields{
		"method": "zipStatusHandler",
	}
	logger.Log(log.InfoLevel, logField, "Zip status handler begins")
	response := Response{}
	folder := r.URL.Query().Get("folder")

	if folder == "" {
		http.Error(w, "Folder not specified", http.StatusBadRequest)
		response.StatusCode = http.StatusBadRequest
		response.Status = "Bad Request"
		response.Message = "Folder not specified"
		response.ResponseTime = time.Now()
		json.NewEncoder(w).Encode(response)
		return
	}
	logger.Log(log.TraceLevel, logField, "Loading zipping status from map")
	status, ok := zipStatuses.Load(folder)
	if !ok {
		http.Error(w, "No zipping process found for the specified folder", http.StatusNotFound)
		response.StatusCode = http.StatusNotFound
		response.Status = "Not Found"
		response.Message = "No zipping process found for the specified folder"
		response.ResponseTime = time.Now()
		json.NewEncoder(w).Encode(response)
		return
	}

	response.StatusCode = http.StatusOK
	response.Status = "OK"
	response.Message = "Zipping completed for the specified folder"
	response.ResponseTime = time.Now()
	json.NewEncoder(w).Encode(response)
	json.NewEncoder(w).Encode(status)
	logger.Log(log.InfoLevel, logField, "Exits zip status handler")
}
