package main

import (
	"encoding/json"
	"log"
	"os"
)

func saveFolderMetadata(){
	metadataLock.Lock()
	defer metadataLock.Unlock()

	var folderMetadataList []FolderMetadata
	folderMetadataMap.Range(func(key, value any) bool {
		folderMetadataList = append(folderMetadataList, value.(FolderMetadata))
		return true
	})
	metaDataFile:=uploadDir+".json"
	folder,err:=os.Create(metaDataFile)
	if err!=nil{
		log.Fatal("Error in openning folder metadata file",err)
		return
	}
	defer folder.Close()
	err=json.NewEncoder(folder).Encode(folderMetadataList)
	if err!=nil{
		log.Fatalln("Error in storing folder metadata",err)
		return
	}
	log.Println("Folder Metadata saved successfully")
}