package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type FolderDetails struct {
	FileCount    int       `json:"file_count"`
	FolderSize   int       `json:"folder_size"`
	CreatedTime  time.Time `json:"created_time"`
	ModifiedTime time.Time `json:"modified_time"`
}

func printStat(file fs.FileInfo, w http.ResponseWriter) {
	fmt.Fprintf(w, "%-30v", file.Mode())
	fmt.Fprintf(w, "%-30v", file.Name())
	fmt.Fprintf(w, "%-30v", file.Size())
	fmt.Fprintf(w, "%-30v\n", file.ModTime().Format(http.TimeFormat))
}

func folderInfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid HTTP request", http.StatusBadRequest)
		return
	}

	folder := r.URL.Query().Get("folder")
	if folder == "" {
		http.Error(w, "Folder not specified", http.StatusBadRequest)
		return
	}
	folderPath := filepath.Join(uploadDir, folder)
	folderInfo, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		log.Fatal(err)
		return
	}
	content, err := os.ReadDir(folderPath)
	if err != nil {
		log.Fatal(err)
		return
	}
	var result any
	metaDataMap.Range(func(key, value any) bool {
		if key == folderPath{
			result=value
		}
		return true;
	})
	fmt.Fprintf(w, "Name: %s\n", folderInfo.Name())
	fmt.Fprintf(w, "Size: %d\n", folderInfo.Size())
	fmt.Fprintf(w, "Count: %d\n", len(content))
	fmt.Fprintf(w, "Last-Modified: %v\n", folderInfo.ModTime().Format(http.TimeFormat))
	filePath := filepath.Join(folderPath, content[0].Name())
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Fprintf(w, "Created-Time: %v\n", fileInfo.ModTime().Format(http.TimeFormat))
	fmt.Fprintf(w, "%-30v", "Permission")
	fmt.Fprintf(w, "%-30v", "Name")
	fmt.Fprintf(w, "%-30v", "Size")
	fmt.Fprintf(w, "%-30v\n", "Modified-time")
	for _, f := range content {
		//fmt.Fprintf(w,"%s\n",f.Name())
		filePath := filepath.Join(folderPath, f.Name())
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			log.Fatal(err)
			return
		}
		printStat(fileInfo, w)
	}
	json.NewEncoder(w).Encode(result)
}
