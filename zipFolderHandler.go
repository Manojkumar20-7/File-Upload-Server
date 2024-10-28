package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func queueZipping(folder string) {
	logger.Log(log.InfoLevel, log.Fields{"method": "worker"}, "Adding zip folder to queue")
	wg.Add(1)
	taskQueue <- folder
}

func workerPool() {
	logger.Log(log.InfoLevel, log.Fields{"method": "worker"}, "Initiating zipping worker pool")
	for i := 0; i < workerCount; i++ {
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
	logField := log.Fields{
		"method": "zipFolderHandler",
	}
	logger.Log(log.InfoLevel, logField, "Zip folder handler begins")
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
	logger.Log(log.DebugLevel, logField, "Initiates zipping folder")
	queueZipping(folder)
	logger.Log(log.TraceLevel, logField, fmt.Sprintf("Zipping process started for fodler %s", folder))
	response.StatusCode = http.StatusOK
	response.Status = "OK"
	response.Message = "Zipping process started for folder " + folder
	response.ResponseTime = time.Now()
	jsonResponse.Encode(response)
	logger.Log(log.InfoLevel, logField, "Zip folder handler exits")
}
