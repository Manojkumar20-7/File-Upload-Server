package main

import (
	"context"
	"fileServer/constants"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	log "fileServer/logger"

	"github.com/sirupsen/logrus"
)

var (
	fileLocks         sync.Map
	metaDataMap       sync.Map
	metadataLock      sync.Mutex
	folderMetadataMap sync.Map
	taskQueue         = make(chan string, 150)
	loadQueue         = make(chan string, 200)
	wg                sync.WaitGroup
	loadwg            sync.WaitGroup
	zipStatuses       sync.Map
	logger            *log.BufferedLogger
)

func init() {
	_, err := os.Stat(constants.UploadDir)
	if os.IsNotExist(err) {
		os.Mkdir(constants.UploadDir, os.ModePerm)
	}

	logger = log.NewBufferedLogger(filepath.Join(constants.LogDir, constants.LogFile), 1, 1, time.Millisecond)
}

func main() {
	logField := logrus.Fields{
		"method": "main",
	}
	logger.Log(logrus.InfoLevel, logField, "Application starts...")
	loadWorkerPool()
	loadFileMetadataAtStart()
	loadFolderMetadata()
	workerPool()
	server := http.Server{
		Addr: ":8080",
	}
	http.HandleFunc("/upload", uploadFileHandler)
	http.HandleFunc("/download", downloadFileHandler)
	http.HandleFunc("/fileinfo", fileInfoHandler)
	http.HandleFunc("/delete", deleteFileHandler)
	http.HandleFunc("/folderinfo", folderInfoHandler)
	http.HandleFunc("/createfolder", folderCreateHandler)
	http.HandleFunc("/deletefolder", deleteFolderHandler)
	http.HandleFunc("/zip", zipFolderHandler)
	http.HandleFunc("/zipdownload", downloadZipHandler)
	http.HandleFunc("/zipstatus", zipStatusHandler)
	logger.Log(logrus.InfoLevel, logField, "Server is listening in http://localhost:8080")
	logrus.Println("Server is listening in http://localhost:8080")
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Log(logrus.FatalLevel, logField, fmt.Sprintf("HTTP shutdown error: %v", err))
		}
		logger.Log(logrus.InfoLevel, logField, "Server is shutting down")
	}()
	err := server.ListenAndServe()
	if err == http.ErrServerClosed {
		fmt.Println("Server is shutting down")
		Shutdown()
	}
}

func Shutdown() {
	close(taskQueue)
	wg.Wait()
	logger.Log(logrus.InfoLevel, logrus.Fields{"method": "Shutdown"}, "All task completed, Shutting down...")
	logrus.Println("All task completed, shutting down...")
}