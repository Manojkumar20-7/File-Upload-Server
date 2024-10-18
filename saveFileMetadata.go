package main

import (
	"encoding/json"
	"log"
	"os"
)

func saveFileMetadata() {
	metadataLock.Lock()
	defer metadataLock.Unlock()

	var metadataList []FileMetadata
	metaDataMap.Range(func(key, value any) bool {
		metadataList = append(metadataList, value.(FileMetadata))
		return true
	})

	file, err := os.Create(fileMetadataFile)
	if err != nil {
		log.Fatalln("Error in opening metadata file while saving", err)
		return
	}
	defer file.Close()
	err = json.NewEncoder(file).Encode(metadataList)
	if err != nil {
		log.Fatalln("Error in storing metadata in file", err)
		return
	}
	log.Println("Metadata saved successfully")
}
