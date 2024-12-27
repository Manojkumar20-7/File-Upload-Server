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

func folderInfoHandler(w http.ResponseWriter, r *http.Request) {
	metrics.ActiveRequest.Inc()
	now:=time.Now()
	metrics.RequestCount.With(prometheus.Labels{
		"path":r.URL.Path,
		"method":r.Method,
	}).Inc()
	logField := log.Fields{
		"method": "folderInfoHandler",
	}
	logger.Log(log.InfoLevel, logField, "Folder info handler begins")
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
	if folder == "" {
		http.Error(w, "Folder not specified", http.StatusBadRequest)
		response.StatusCode = http.StatusBadRequest
		response.Status = "Bad Request"
		response.Message = "Folder not specified"
		response.ResponseTime = time.Now()
		jsonResponse.Encode(response)
		return
	}
	folderPath := filepath.Join(constants.UploadDir, folder)
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
	go getAllMetrics()
	response.StatusCode = http.StatusOK
	response.Status = "OK"
	response.Message = "Folder info retrieved successfully"
	response.ResponseTime = time.Now()
	jsonResponse.Encode(response)
	json.NewEncoder(w).Encode(result)
	logger.Log(log.InfoLevel, logField, "Folder info handler exits")
	time.Sleep(time.Second*time.Duration(rand.Intn(15)))
	metrics.ResponseTime.Observe(float64(time.Since(now).Seconds()))
	metrics.RequestTime.With(prometheus.Labels{"path":r.URL.Path}).Observe(float64(time.Since(now)))
	metrics.ActiveRequest.Dec()
}
