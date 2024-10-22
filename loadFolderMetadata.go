package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

func loadFolderMetadata() {
	metadataFile:=uploadDir+".json"
	folder, err := os.Open(metadataFile)
	if err != nil {
		log.Fatal("Error in opening metadata file")
		return
	}
	defer folder.Close()

	var metadataList []FolderMetadata
	if err := json.NewDecoder(folder).Decode(&metadataList); err != nil && err!=io.EOF {
		log.Fatal("Failed to decode the metadata", err)
	}

	for _, metadata := range metadataList {
		folderMetadataMap.Store(metadata.FolderPath, metadata)
		log.Println("Loaded: ", metadata.FolderPath)
	}

	log.Println("Folder Metadata loaded successfully in memory")
}
