package main

import (
	"encoding/json"
	"fileServer/config"
	"fileServer/constants"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

func queueZipping(folder string) {
	logger.Log(log.InfoLevel, log.Fields{"method": "worker"}, "Adding zip folder to queue")
	wg.Add(1)
	taskQueue <- folder
}

func workerPool() {
	logger.Log(log.InfoLevel, log.Fields{"method": "worker"}, "Initiating zipping worker pool")
	for i := 0; i < constants.WorkerCount; i++ {
		go worker(i)
	}
}

func worker(workerID int) {
	logger.Log(log.TraceLevel, log.Fields{"method": "worker"}, "Starts zipping folder in background")
	for folder := range taskQueue {
		logger.Log(log.TraceLevel, log.Fields{"method": "worker"}, fmt.Sprintf("Worker %d: Started zipping folder: %s\n", workerID, folder))
		zipFolderInBackground(folder)
		wg.Done()
		logger.Log(log.TraceLevel, log.Fields{"method": "worker"}, fmt.Sprintf("Worker %d: Finished Zipping folder %s\n", workerID, folder))
	}
}

func zipFolderHandler(w http.ResponseWriter, r *http.Request) {
	now:=time.Now()
	metrics.ActiveRequest.Inc()
	metrics.RequestCount.With(prometheus.Labels{
		"path":r.URL.Path,
		"method":r.Method,
	}).Inc()
	logField := log.Fields{
		"method": "zipFolderHandler",
	}
	logger.Log(log.InfoLevel, logField, "Zip folder handler begins")
	response := config.Response{}
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
	go getAllMetrics()
	logger.Log(log.DebugLevel, logField, "Initiates zipping folder")
	queueZipping(folder)
	logger.Log(log.TraceLevel, logField, fmt.Sprintf("Zipping process started for fodler %s", folder))
	response.StatusCode = http.StatusOK
	response.Status = "OK"
	response.Message = "Zipping process started for folder " + folder
	response.ResponseTime = time.Now()
	jsonResponse.Encode(response)
	logger.Log(log.InfoLevel, logField, "Zip folder handler exits")
	time.Sleep(time.Second*time.Duration(rand.Intn(15)))
	metrics.ResponseTime.Observe(float64(time.Since(now).Seconds()))
	metrics.RequestTime.With(prometheus.Labels{"path":r.URL.Path}).Observe(float64(time.Since(now)))
	metrics.ActiveRequest.Dec()
}
