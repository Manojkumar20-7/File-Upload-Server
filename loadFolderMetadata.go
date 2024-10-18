package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

func loadFolderMetadata() {
	file, err := os.Open(folderMetadataFile)
	if err != nil {
		log.Fatal("Error in opening metadata file")
		return
	}
	defer file.Close()

	var metadataList []FolderMetadata
	if err := json.NewDecoder(file).Decode(&metadataList); err != nil && err!=io.EOF {
		log.Fatal("Failed to decode the metadata", err)
	}

	for _, metadata := range metadataList {
		folderMetadataMap.Store(metadata.FolderPath, metadata)
		log.Println("Loaded: ", metadata.FolderPath)
	}

	log.Println("Folder Metadata loaded successfully in memory")
}
