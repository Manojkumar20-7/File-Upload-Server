package main

import (
	"io/fs"
	"log"
	"path/filepath"
)

func loadFileMetadataAtStart()  {
	err:=filepath.WalkDir(uploadDir,func(path string, d fs.DirEntry, err error) error {
		if d.IsDir(){
			loadFileMetadata(path)
			return nil
		}
		return err
	})
	if err!=nil{
		log.Fatalln("Error in loadind metadata at startup", err)
		return
	}
}