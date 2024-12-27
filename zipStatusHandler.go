package main

import (
	"encoding/json"
	"fileServer/config"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

func zipStatusHandler(w http.ResponseWriter, r *http.Request) {
	metrics.ActiveRequest.Inc()
	now:=time.Now()
	metrics.RequestCount.With(prometheus.Labels{
		"path":r.URL.Path,
		"method":r.Method,
	}).Inc()
	logField := log.Fields{
		"method": "zipStatusHandler",
	}
	logger.Log(log.InfoLevel, logField, "Zip status handler begins")
	response := config.Response{}
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
	time.Sleep(time.Second*time.Duration(rand.Intn(15)))
	metrics.ResponseTime.Observe(float64(time.Since(now).Seconds()))
	metrics.RequestTime.With(prometheus.Labels{"path":r.URL.Path}).Observe(float64(time.Since(now)))
	metrics.ActiveRequest.Dec()
}
