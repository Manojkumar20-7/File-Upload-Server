package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

func loadFileMetadata(folderPath string) {
	logField:=log.Fields{
		"method":"loadFileMetadata",
	}
	logger.Log(log.TraceLevel,logField,"Load file metadata started")
	metadataFile := folderPath + ".json"
	file, err := os.OpenFile(metadataFile, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		logger.Log(log.FatalLevel,logField,"Error in opening metadata file")
		return
	}
	defer file.Close()

	var metadataList []FileMetadata
	if err := json.NewDecoder(file).Decode(&metadataList); err != nil && err != io.EOF {
		logger.Log(log.FatalLevel,logField,"Failed to decode the metadata")
	}

	for _, metadata := range metadataList {
		metaDataMap.Store(metadata.FilePath, metadata)
		logger.Log(log.TraceLevel,logField,fmt.Sprintf("Loaded file: %s", metadata.FilePath))
	}

	logger.Log(log.TraceLevel,logField,"Metadata loaded successfully in memory")
}
