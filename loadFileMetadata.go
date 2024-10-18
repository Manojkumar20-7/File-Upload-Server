package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

func loadFileMetadata() {
	file, err := os.Open(fileMetadataFile)
	if err != nil {
		log.Fatal("Error in opening metadata file")
		return
	}
	defer file.Close()

	var metadataList []FileMetadata
	if err := json.NewDecoder(file).Decode(&metadataList); err != nil && err!=io.EOF{
		log.Fatal("Failed to decode the metadata", err)
	}

	for _, metadata := range metadataList {
		metaDataMap.Store(metadata.FilePath, metadata)
		log.Println("Loaded: ", metadata.FilePath)
	}

	log.Println("Metadata loaded successfully in memory")
}
