package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{}
	jsonResponse := json.NewEncoder(w)
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid HTTP method", http.StatusBadRequest)
		response.StatusCode = http.StatusBadRequest
		response.Status = "Bad Resquest"
		response.Message = "Invalid HTTP method"
		response.ResponseTime = time.Now()
		jsonResponse.Encode(response)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Input not specified", http.StatusBadRequest)
		response.StatusCode = http.StatusBadRequest
		response.Status = "Bad Resquest"
		response.Message = "Input not specified"
		response.ResponseTime = time.Now()
		jsonResponse.Encode(response)
		return
	}

	folder := r.FormValue("folder")
	fileName := r.FormValue("filename")

	if folder == "" {
		http.Error(w, "Folder not specified", http.StatusBadRequest)
		response.StatusCode = http.StatusBadRequest
		response.Status = "Bad Resquest"
		response.Message = "Folder not specified"
		response.ResponseTime = time.Now()
		jsonResponse.Encode(response)
		return
	}
	if fileName == "" {
		http.Error(w, "Filename not specified", http.StatusBadRequest)
		response.StatusCode = http.StatusBadRequest
		response.Status = "Bad Resquest"
		response.Message = "Filename not specified"
		response.ResponseTime = time.Now()
		jsonResponse.Encode(response)
		return
	}

	f, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File not specified", http.StatusBadRequest)
		response.StatusCode = http.StatusBadRequest
		response.Status = "Bad Resquest"
		response.Message = "File not specified"
		response.ResponseTime = time.Now()
		jsonResponse.Encode(response)
		return
	}
	folderPath := filepath.Join(uploadDir, folder)

	folderLock := getFileLock(folderPath)
	folderLock.Lock()
	defer folderLock.Unlock()

	_, er := os.Stat(folderPath)
	if os.IsNotExist(er) {
		os.MkdirAll(folderPath, os.ModePerm)
	}
	folderInfo, err := os.Stat(folderPath)
	if err != nil {
		http.Error(w, "Unable to stat folder", http.StatusInternalServerError)
		response.StatusCode = http.StatusInternalServerError
		response.Status = "Internal Server Error"
		response.Message = "Unable to stat folder"
		response.ResponseTime = time.Now()
		jsonResponse.Encode(response)
		return
	}
	
	filePath := filepath.Join(folderPath, fileName)

	fileLock := getFileLock(filePath)
	fileLock.Lock()
	defer fileLock.Unlock()

	file, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Unable to create file", http.StatusInternalServerError)
		response.StatusCode = http.StatusInternalServerError
		response.Status = "Internal Server Error"
		response.Message = "Unable to create file"
		response.ResponseTime = time.Now()
		jsonResponse.Encode(response)
		return
	}
	defer file.Close()
	_, err = io.Copy(file, f)
	if err != nil {
		http.Error(w, "Error in writing file", http.StatusInternalServerError)
		response.StatusCode = http.StatusInternalServerError
		response.Status = "Internal Server Error"
		response.Message = "Error in writing file"
		response.ResponseTime = time.Now()
		jsonResponse.Encode(response)
		return
	}

	fileInfo, _ := os.Stat(filePath)

	metadata := FileMetadata{
		FileName:     fileName,
		FilePath:     filePath,
		FolderPath:   folderPath,
		FileSize:     fileInfo.Size(),
		ModifiedTime: fileInfo.ModTime().Format(http.TimeFormat),
		CreatedTime:  time.Now().Format(http.TimeFormat),
		FileMode:     fileInfo.Mode(),
		IsDirectory:  fileInfo.IsDir(),
	}
	metaDataMap.Store(filePath, metadata)
	saveFileMetadata(folderPath)
	folderMetadata := FolderMetadata{FilesCount: 0}
	folderMetadata.FolderName = folderInfo.Name()
	folderMetadata.FolderPath = folderPath
	folderMetadata.FolderSize = folderInfo.Size()
	folderMetadata.FilesCount = folderMetadata.FilesCount+1
	folderMetadata.ModifiedTime = folderInfo.ModTime().Format(http.TimeFormat)
	folderMetadata.CreatedTime = time.Now().Format(http.TimeFormat)
	folderMetadata.FolderMode = folderInfo.Mode()
	folderMetadata.IsDirectory = folderInfo.IsDir()
	folderMetadataMap.Store(folderPath, folderMetadata)
	saveFolderMetadata()
	response.StatusCode = http.StatusCreated
	response.Status = "Created"
	response.Message = "File uploaded successfully"
	response.ResponseTime = time.Now()
	jsonResponse.Encode(response)
}
