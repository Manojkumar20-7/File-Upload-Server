package main

import (
	"encoding/json"
	"fileServer/config"
	"os"

	"github.com/sirupsen/logrus"
)

func saveFileMetadata(folderPath string) {
	logField := logrus.Fields{
		"method": "saveFileMetadata",
	}
	logger.Log(logrus.InfoLevel, logField, "Save file metadata begin")
	metadataLock.Lock()
	logger.Log(logrus.DebugLevel, logField, "Acquiring lock for map")
	defer metadataLock.Unlock()

	var metadataList []config.FileMetadata
	metaDataMap.Range(func(key, value any) bool {
		if value.(config.FileMetadata).FolderPath == folderPath {
			metadataList = append(metadataList, value.(config.FileMetadata))
		}
		return true
	})

	metadataFile := folderPath + ".json"
	file, err := os.Create(metadataFile)
	if err != nil {
		logger.Log(logrus.ErrorLevel, logField, "Error in opening metadata file while saving")
		return
	}
	defer file.Close()
	logger.Log(logrus.DebugLevel, logField, "Saving metadata into json")
	err = json.NewEncoder(file).Encode(metadataList)
	if err != nil {
		logger.Log(logrus.FatalLevel, logField, "Error in storing metadata in file")
		return
	}
	logger.Log(logrus.TraceLevel, logField, "Metadata saved successfully")
	logger.Log(logrus.InfoLevel, logField, "Exits save file metada")
}
