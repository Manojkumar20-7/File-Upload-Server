package main

import (
	"fmt"
	"io/fs"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func queueLoading(folderPath string) {
	loadwg.Add(1)
	loadQueue <- folderPath
}

func loadWorkerPool() {
	logger.Log(log.DebugLevel, log.Fields{"method": "loadWorkerPool"}, "Initiating load worker pool")
	for i := 0; i < loadWorkerCount; i++ {
		go loadWorker(i)
	}
}

func loadWorker(workerID int) {
	logger.Log(log.TraceLevel, log.Fields{"method": "loadWorker"}, "Initiating load worker pool")
	for folder := range loadQueue {
		logger.Log(log.TraceLevel, log.Fields{"method": "loadWorker"},fmt.Sprintf("Loading worker %d: Started loading metadata: %s\n", workerID, folder))
		loadFileMetadata(folder)
		loadwg.Done()
		logger.Log(log.TraceLevel, log.Fields{"method": "loadWorker"},fmt.Sprintf("Loading worker %d: Finished loading metadata: %s\n", workerID, folder))
	}
	logger.Log(log.DebugLevel, log.Fields{"method": "loadWorker"}, "Loading file metadata into memory is finished")
}

func loadFileMetadataAtStart() {
	logField := log.Fields{
		"method": "loadFileMetadataAtStart",
	}
	logger.Log(log.InfoLevel, logField, "LoadFileMetadata file has initiated")
	err := filepath.WalkDir(uploadDir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			queueLoading(path)
			return nil
		}
		return err
	})
	if err != nil {
		logger.Log(log.FatalLevel, logField, "Error in loading file meta at start")
		return
	}
}
