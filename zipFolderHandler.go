package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func queueZipping(folder string) {
	wg.Add(1)
	taskQueue <- folder
}

func workerPool() {
	for i := 0; i < workerCount; i++ {
		go worker(i)
	}
}

func worker(workerID int) {
	for folder := range taskQueue {
		log.Printf("Worker %d: Started zipping folder: %s\n", workerID, folder)
		zipFolderInBackground(folder)
		wg.Done()
		log.Printf("Worker %d: Finished Zipping folder %s\n", workerID, folder)
	}
}

func zipFolderHandler(w http.ResponseWriter, r *http.Request) {
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

	queueZipping(folder)
	fmt.Fprintf(w, "Zipping process started for fodler %s\n", folder)
	response.StatusCode = http.StatusOK
	response.Status = "OK"
	response.Message = "Zipping process started for folder " + folder
	response.ResponseTime = time.Now()
	jsonResponse.Encode(response)
}
