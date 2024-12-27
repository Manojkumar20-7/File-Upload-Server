package main

import (
	"encoding/json"
	"fileServer/config"
	"fileServer/constants"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

func downloadZipHandler(w http.ResponseWriter, r *http.Request) {
	metrics.ActiveRequest.Inc()
	now:=time.Now()
	metrics.RequestCount.With(prometheus.Labels{
		"path":r.URL.Path,
		"method":r.Method,
	}).Inc()
	logField := log.Fields{
		"method": "downloadZipHandler",
	}
	logger.Log(log.InfoLevel, logField, "Download zip handler begins")
	response := config.Response{}
	jsonResponse := json.NewEncoder(w)
	folder := r.URL.Query().Get("folder")
	if folder == "" {
		logger.Log(log.InfoLevel, logField, "Folder not specified")
		http.Error(w, "Folder not specified", http.StatusBadRequest)
		response.StatusCode = http.StatusBadRequest
		response.Status = "Bad Request"
		response.Message = "Folder not specified"
		response.ResponseTime = time.Now()
		jsonResponse.Encode(response)
		return
	}

	zipPath := filepath.Join(constants.UploadDir, folder)
	zipName := zipPath + ".zip"
	logger.Log(log.DebugLevel, logField, "Checking zip details")
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

	logger.Log(log.DebugLevel, logField, "Loading zip status from map")
	status, ok := zipStatuses.Load(folder)
	if !ok {
		http.Error(w, "No zipping found for the specified folder", http.StatusNotFound)
		response.StatusCode = http.StatusNotFound
		response.Status = "Not Found"
		response.Message = "No zipping found for the specified folder"
		response.ResponseTime = time.Now()
		jsonResponse.Encode(response)
		return
	} else if status.(*config.ZipStatus).Status != "completed" {
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
		zipFilePath := status.(*config.ZipStatus).FilePath
		http.ServeFile(w, r, zipFilePath)
		log.Println(response)
		jsonResponse.Encode(response)
	}
	go getAllMetrics()
	logger.Log(log.InfoLevel, logField, "Exits Zip  download handler")
	time.Sleep(time.Second*time.Duration(rand.Intn(15)))
	metrics.ResponseTime.Observe(float64(time.Since(now).Seconds()))
	metrics.RequestTime.With(prometheus.Labels{"path":r.URL.Path}).Observe(float64(time.Since(now)))
	metrics.ActiveRequest.Dec()
}
