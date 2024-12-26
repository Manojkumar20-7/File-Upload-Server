package main

import (
	"encoding/json"
	"fileServer/config"
	"fileServer/constants"
	"os"

	log "github.com/sirupsen/logrus"
)

func saveFolderMetadata() {
	logField := log.Fields{
		"method": "saveFolderMetadata",
	}
	logger.Log(log.InfoLevel, logField, "Save folder metadata begins")
	metadataLock.Lock()
	defer metadataLock.Unlock()

	var folderMetadataList []config.FolderMetadata
	folderMetadataMap.Range(func(key, value any) bool {
		folderMetadataList = append(folderMetadataList, value.(config.FolderMetadata))
		return true
	})
	metaDataFile := constants.UploadDir + ".json"
	folder, err := os.Create(metaDataFile)
	if err != nil {
		logger.Log(log.ErrorLevel, logField, "Error in openning folder metadata file")
		log.Fatal("Error in openning folder metadata file", err)
		return
	}
	defer folder.Close()
	err = json.NewEncoder(folder).Encode(folderMetadataList)
	if err != nil {
		logger.Log(log.FatalLevel, logField, "Error in saving metadata into file")
		return
	}
	logger.Log(log.TraceLevel, logField, "Folder Metadata saved successfully")
	logger.Log(log.InfoLevel, logField, "Exits folder metadata")
}
