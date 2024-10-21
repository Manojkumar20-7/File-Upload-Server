package main

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func deleteFolderHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{}
	jsonResponse := json.NewEncoder(w)
	folder := r.URL.Query().Get("folder")
	folderPath := filepath.Join(uploadDir, folder)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		folderMetadataMap.Delete(folderPath)
		response.StatusCode = http.StatusOK
		response.Status = "OK"
		response.Message = "Folder deleted successfully"
		response.ResponseTime = time.Now()
		jsonResponse.Encode(response)
		os.Remove(folderPath+".json")
		return
	} else {
		os.RemoveAll(folderPath)
		response.StatusCode = http.StatusOK
		response.Status = "OK"
		response.Message = "Folder " + folder + " deleted"
		response.ResponseTime = time.Now()
		jsonResponse.Encode(response)
		os.Remove(folderPath+".json")
		return
	}
}
