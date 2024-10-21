package main

import (
	"encoding/json"
	"log"
	"os"
)

func saveFileMetadata(folderPath string) {
	metadataLock.Lock()
	defer metadataLock.Unlock()

	var metadataList []FileMetadata
	metaDataMap.Range(func(key , value any) bool {
		if value.(FileMetadata).FolderPath==folderPath{
			metadataList = append(metadataList, value.(FileMetadata))
		}
		return true
	})

	metadataFile:= folderPath+".json"
	file, err := os.Create(metadataFile)
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
