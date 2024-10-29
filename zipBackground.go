package main

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func zipFolderInBackground(folder string) {
	logField := log.Fields{
		"method": "zipFolderInBackground",
	}
	logger.Log(log.InfoLevel, logField, "Background zip begins")
	status := &zipStatus{
		Status:    "in_progress",
		StartTime: time.Now(),
	}
	zipStatuses.Store(folder, status)

	zipName := folder + ".zip"
	zipPath := filepath.Join(uploadDir, zipName)

	zipFile, err := os.Create(zipPath)
	if err != nil {
		logger.Log(log.ErrorLevel, logField, "Unable to create zip file")
		status.Status = "failed"
		status.ErrorMsg = "Unable to create zip file"
		zipStatuses.Store(folder, status)
		return
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	folderPath := filepath.Join(uploadDir, folder)
	logger.Log(log.TraceLevel, logField, "Copying content for zipping")
	err = filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = strings.TrimPrefix(path, folderPath+"/")
		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		fileContent, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		_, err = writer.Write(fileContent)
		return err
	})

	if err != nil {
		logger.Log(log.ErrorLevel, logField, "Error in zipping")
		status.Status = "failed"
		status.ErrorMsg = "Error zipping folder"
	} else {
		logger.Log(log.DebugLevel, logField, "Zipping completed")
		status.Status = "completed"
		status.FilePath = zipPath
	}
	status.EndTime = time.Now()
	zipStatuses.Store(folder, status)
	logger.Log(log.InfoLevel, logField, "Exits Zip in background method")
}
