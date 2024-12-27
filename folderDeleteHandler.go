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

func deleteFolderHandler(w http.ResponseWriter, r *http.Request) {
	metrics.ActiveRequest.Inc()
	now:=time.Now()
	metrics.RequestCount.With(prometheus.Labels{
		"path":r.URL.Path,
		"method":r.Method,
	}).Inc()
	logField := log.Fields{
		"method": "deleteFolderHandler",
	}
	logger.Log(log.InfoLevel, logField, "Delete folder handler begins")
	response := config.Response{}
	jsonResponse := json.NewEncoder(w)
	folder := r.URL.Query().Get("folder")
	folderPath := filepath.Join(constants.UploadDir, folder)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		logger.Log(log.DebugLevel, logField, "Folder not found and updating folder metadata")
		folderMetadataMap.Delete(folderPath)
		saveFolderMetadata()
		response.StatusCode = http.StatusOK
		response.Status = "OK"
		response.Message = "Folder deleted successfully"
		response.ResponseTime = time.Now()
		jsonResponse.Encode(response)
		os.Remove(folderPath + ".json")
	} else {
		logger.Log(log.TraceLevel, logField, "Removing folder")
		os.RemoveAll(folderPath)
		logger.Log(log.TraceLevel, logField, "Updating folder metadata")
		folderMetadataMap.Delete(folderPath)
		response.StatusCode = http.StatusOK
		response.Status = "OK"
		response.Message = "Folder " + folder + " deleted"
		response.ResponseTime = time.Now()
		jsonResponse.Encode(response)
		os.Remove(folderPath + ".json")
		saveFolderMetadata()
	}
	go getAllMetrics()
	logger.Log(log.InfoLevel, logField, "Delete folder handler exits")
	time.Sleep(time.Second*time.Duration(rand.Intn(15)))
	metrics.ResponseTime.Observe(float64(time.Since(now).Seconds()))
	metrics.RequestTime.With(prometheus.Labels{"path":r.URL.Path}).Observe(float64(time.Since(now)))
	metrics.ActiveRequest.Dec()
}
