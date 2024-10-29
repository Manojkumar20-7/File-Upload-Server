package main

import (
	"encoding/json"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
)

func getCurrentTime(folderPath string) string {
	var result string
	folderMetadataMap.Range(func(key, value any) bool {
		if folderPath == key && value.(FolderMetadata).CreatedTime != "" {
			result = value.(FolderMetadata).CreatedTime
		}
		return true
	})
	return result
}

func getFilesCount(folderPath string) (int, error) {
	fcount := 0
	_ = filepath.WalkDir(folderPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		fcount++
		return nil
	})
	return fcount, nil
}

func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	logField := log.Fields{
		"method": "uploadFileHandler",
	}
	logger.Log(log.InfoLevel, logField, "Upload handler begins")
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
	logger.Log(log.DebugLevel, logField, "Reading request body at upload handler")
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

	logger.Log(log.DebugLevel, logField, "Acquiring folder lock at upload handler")
	folderLock := getFileLock(folderPath)
	folderLock.Lock()
	defer folderLock.Unlock()

	_, er := os.Stat(folderPath)
	if os.IsNotExist(er) {
		logger.Log(log.TraceLevel, logField, "Creating folder")
		os.MkdirAll(folderPath, os.ModePerm)
	}

	filePath := filepath.Join(folderPath, fileName)
	logger.Log(log.DebugLevel, logField, "Acquiring file lock at upload handler")
	fileLock := getFileLock(filePath)
	fileLock.Lock()
	defer fileLock.Unlock()

	file, err := os.Create(filePath)
	if err != nil {
		logger.Log(log.ErrorLevel, logField, "Error in creating file")
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
		logger.Log(log.ErrorLevel, logField, "Error in copying file content")
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
	logger.Log(log.DebugLevel, logField, "Storing file metadata in memory")
	metaDataMap.Store(filePath, metadata)
	saveFileMetadata(folderPath)
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
	folderMetadata := FolderMetadata{}
	folderMetadata.FolderName = folderInfo.Name()
	folderMetadata.FolderPath = folderPath
	folderMetadata.FolderSize = folderInfo.Size()
	fcount, _ := getFilesCount(folderPath)
	folderMetadata.FilesCount = fcount
	folderMetadata.ModifiedTime = folderInfo.ModTime().Format(http.TimeFormat)
	curTime := getCurrentTime(folderPath)
	if curTime == "" {
		folderMetadata.CreatedTime = time.Now().Format(http.TimeFormat)
	} else {
		folderMetadata.CreatedTime = curTime
	}
	folderMetadata.FolderMode = folderInfo.Mode().Perm()
	folderMetadata.IsDirectory = folderInfo.IsDir()
	folderMetadataMap.Store(folderPath, folderMetadata)
	saveFolderMetadata()
	response.StatusCode = http.StatusCreated
	response.Status = "Created"
	response.Message = "File uploaded successfully"
	response.ResponseTime = time.Now()
	jsonResponse.Encode(response)
	logger.Log(log.InfoLevel, logField, "Completing file uploads and exits")
}
