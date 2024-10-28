package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

func loadFolderMetadata() {
	logField := log.Fields{
		"method": "loadFolderMetadata",
	}
	logger.Log(log.InfoLevel, logField, "Loading Folder metadata starts")
	metadataFile := uploadDir + ".json"
	folder, err := os.Open(metadataFile)
	if err != nil {
		logger.Log(log.FatalLevel, logField, "Error in opening metadata file")
		return
	}
	defer folder.Close()

	var metadataList []FolderMetadata
	if err := json.NewDecoder(folder).Decode(&metadataList); err != nil && err != io.EOF {
		logger.Log(log.FatalLevel, logField, "Failed to decode the metadata")
	}

	for _, metadata := range metadataList {
		folderMetadataMap.Store(metadata.FolderPath, metadata)
		if os.IsNotExist(err) {
			logger.Log(log.FatalLevel, logField, "Error in Loading Folder metadata into map")
			return
		}
		logger.Log(log.TraceLevel, logField, fmt.Sprintf("Loaded folder: %s", metadata.FolderPath))
	}

	logger.Log(log.DebugLevel, logField, "Folder Metadata loaded successfully in memory")
}
