package main

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
)

func deleteFolderHandler(w http.ResponseWriter, r *http.Request) {
	logField := log.Fields{
		"method": "deleteFolderHandler",
	}
	logger.Log(log.InfoLevel, logField, "Delete folder handler begins")
	response := Response{}
	jsonResponse := json.NewEncoder(w)
	folder := r.URL.Query().Get("folder")
	folderPath := filepath.Join(uploadDir, folder)
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
	logger.Log(log.InfoLevel, logField, "Delete folder handler exits")
}
