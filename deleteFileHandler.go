package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func deleteFileHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{}
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid HTTP request", http.StatusBadRequest)
		response.StatusCode = http.StatusBadRequest
		response.Status = "Bad Request"
		response.Message = "Invalid HTTP resquest"
		response.ResponseTime = time.Now()
		json.NewEncoder(w).Encode(response)
		return
	}

	folder := r.URL.Query().Get("folder")
	if folder == "" {
		http.Error(w, "Folder not specified", http.StatusBadRequest)
		response.StatusCode = http.StatusBadRequest
		response.Status = "Folder not specified"
		response.Message = "Invalid HTTP resquest"
		response.ResponseTime = time.Now()
		json.NewEncoder(w).Encode(response)
		return
	}

	folderPath := filepath.Join(uploadDir, folder)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		http.Error(w, "Folder Not Found", http.StatusNotFound)
		response.StatusCode = http.StatusNotFound
		response.Status = "Not Found"
		response.Message = "Folder Not Found"
		response.ResponseTime = time.Now()
		json.NewEncoder(w).Encode(response)
		return
	}
	if !r.URL.Query().Has("filename") {
		_, err := os.Stat(folderPath)
		if err != nil {
			response.StatusCode = http.StatusBadRequest
			response.Status = "Bad Request"
			response.Message = "Filename not specified"
			response.ResponseTime = time.Now()
			json.NewEncoder(w).Encode(response)
			return
		} else {
			err := os.RemoveAll(folderPath)
			os.Remove(folderPath + ".json")
			if err != nil {
				log.Fatal("Error in deleting folder")
				response.StatusCode = http.StatusBadRequest
				response.Status = "Bad Request"
				response.Message = "File not specified"
				response.ResponseTime = time.Now()
				json.NewEncoder(w).Encode(response)
				return
			}
			folderMetadataMap.Delete(folderPath)
			saveFolderMetadata()
			fmt.Fprintf(w, "Folder deleted successfully: %s\n", folderPath)
			response.StatusCode = http.StatusOK
			response.Status = "OK"
			response.Message = "Folder deleted successfully"
			response.ResponseTime = time.Now()
			json.NewEncoder(w).Encode(response)
		}
	}
	fileName := r.URL.Query().Get("filename")
	filePath := filepath.Join(folderPath, fileName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		response.StatusCode = http.StatusOK
		response.Status = "OK"
		response.Message = "File deleted successfully"
		response.ResponseTime = time.Now()
		json.NewEncoder(w).Encode(response)
		return
	}
	fileLock := getFileLock(filePath)
	fileLock.Lock()
	defer fileLock.Unlock()

	err := os.Remove(filePath)

	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Status = "Internal Server Error"
		response.Message = "Error in deleting file"
		response.ResponseTime = time.Now()
		json.NewEncoder(w).Encode(response)
		return
	}

	metaDataMap.Delete(filePath)
	saveFileMetadata(folderPath)
	fmt.Fprintf(w, "File deleted successfully: %s\n", fileName)
	response.StatusCode = http.StatusOK
	response.Status = "OK"
	response.Message = "File deleted successfully"
	response.ResponseTime = time.Now()
	json.NewEncoder(w).Encode(response)
	file, err := os.Stat(folderPath)
	if err != nil {
		http.Error(w, "File Not Found", http.StatusNotFound)
		response.StatusCode = http.StatusNotFound
		response.Status = "Not Found"
		response.Message = "File Not Found"
		response.ResponseTime = time.Now()
		json.NewEncoder(w).Encode(response)
		return
	}
	if file.Size() == 0 {
		folderLock := getFileLock(folderPath)
		folderLock.Lock()
		defer folderLock.Unlock()

		err := os.Remove(folderPath)
		if err != nil {
			http.Error(w, "Error in deleting folder", http.StatusInternalServerError)
			response.StatusCode = http.StatusInternalServerError
			response.Status = "Internal Server Error"
			response.Message = "Error in deleting folder"
			response.ResponseTime = time.Now()
			json.NewEncoder(w).Encode(response)
			return
		}
		folderMetadataMap.Delete(folderPath)
		saveFolderMetadata()
		os.Remove(folderPath + ".json")
		response.StatusCode = http.StatusOK
		response.Status = "OK"
		response.Message = "Folder deleted successfully"
		response.ResponseTime = time.Now()
		json.NewEncoder(w).Encode(response)
	}
	folderInfo, err := os.Stat(folderPath)
	if err!=nil|| os.IsNotExist(err){
		return
	}
	folderMetadata := FolderMetadata{}
	folderMetadata.FolderName = folderInfo.Name()
	folderMetadata.FolderPath = folderPath
	folderMetadata.FolderSize = folderInfo.Size()
	fcount,_:=getFilesCount(folderPath)
	folderMetadata.FilesCount = fcount
	folderMetadata.ModifiedTime = folderInfo.ModTime().Format(http.TimeFormat)
	folderMetadata.CreatedTime = time.Now().Format(http.TimeFormat)
	folderMetadata.FolderMode = folderInfo.Mode()
	folderMetadata.IsDirectory = folderInfo.IsDir()
	folderMetadataMap.Store(folderPath, folderMetadata)
	saveFolderMetadata()
}
