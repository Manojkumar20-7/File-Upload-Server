package main

import (
	"encoding/json"
	"log"
	"os"
)

func saveFolderMetadata(){
	metadataLock.Lock()
	defer metadataLock.Unlock()

	var folderMetadataLust []FolderMetadata
	folderMetadataMap.Range(func(key, value any) bool {
		folderMetadataLust = append(folderMetadataLust, value.(FolderMetadata))
		return true
	})

	folder,err:=os.Create(folderMetadataFile)
	if err!=nil{
		log.Fatal("Error in openning folder metadata file",err)
		return
	}
	defer folder.Close()
	err=json.NewEncoder(folder).Encode(folderMetadataLust)
	if err!=nil{
		log.Fatalln("Error in storing folder metadata",err)
		return
	}
	log.Println("Folder Metadata saved successfully")
}