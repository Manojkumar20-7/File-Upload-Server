package main

import (
	"encoding/json"
	"fileServer/config"
	"fileServer/constants"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

func fileInfoHandler(w http.ResponseWriter, r *http.Request) {
	metrics.ActiveRequest.Inc()
	now:=time.Now();
	metrics.RequestCount.With(prometheus.Labels{
		"path":r.URL.Path,
		"method":r.Method,
	}).Inc()
	logField:=log.Fields{
		"method":"fileInfoHandler",
	}
	logger.Log(log.InfoLevel,logField,"FileInfo handler begins")
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
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		response.StatusCode = http.StatusNotFound
		response.Status = "Not Found"
		response.Message = "File not found"
		response.ResponseTime = time.Now()
		jsonResponse.Encode(response)
		return
	}
	var result any
	metaDataMap.Range(func(key, value any) bool {
		if key == filePath {
			result = value
		}
		return true
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	go getAllMetrics()
	response.StatusCode = http.StatusOK
	response.Status = "OK"
	response.Message = "File info retrieved successfully"
	response.ResponseTime = time.Now()
	jsonResponse.Encode(response)
	json.NewEncoder(w).Encode(result)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	w.Header().Set("Last-Modified", fileInfo.ModTime().Format(http.TimeFormat))
	logger.Log(log.InfoLevel,logField,"Completing fileinfo handler and exits")
	time.Sleep(time.Second*time.Duration(rand.Intn(15)))
	metrics.ResponseTime.Observe(float64(time.Since(now).Seconds()))
	metrics.RequestTime.With(prometheus.Labels{"path":r.URL.Path}).Observe(float64(time.Since(now)))
	metrics.ActiveRequest.Dec()
}
