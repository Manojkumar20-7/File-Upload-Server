# Go File Server

This project implements a file server using Golang that supports the following operations:

- Upload files
- Download files
- Retrieve file details
- Delete files
- Create empty folders
- Delete folders
- Retrieve folder details
- Zip folders in the background
- Fetch zipping status
- Download zipped folders

## Installation and setup
### Prerequisites
- [Go 1.18+](https://go.dev/dl/) installed
- Any HTTP client (eg., curl, Postman)
### Clone the repository
```bash
git clone https://github.com/Manojkumar20-7/File-Upload-Server.git
cd File-Upload-Server
```
### Install dependencies
```bash
go mod tidy
```
### Run the server
```bash
go run *.go
```


## Endpoints

### 1. File Upload

Upload a file to a specific folder.
```bash
curl -X POST 'http://localhost:8080/upload' \
--form 'folder="FolderName"' \
--form 'filename="Filename.ext"' \
--form 'file=@"path/of/your/file/to/upload"'
```
### 2. File Download

```bash
curl 'http://localhost:8080/download?folder=FolderName&filename=FileName.ext' -o path/to/save/the/downloaded/file/filename.ext
```
### 3. File Info

```bash
curl 'http://localhost:8080/fileinfo?folder=FolderName&filename=Filename.ext'
```
### 4. File Delete
If you delete the last file in a folder, the folder will also be deleted.

```bash
curl -X DELETE 'http://localhost:8080/delete?folder=FolerName&filename=FileName.ext'
```
### 5. Create Folder
```bash
curl 'http://localhost:8080/createfolder?folder=FolderName'
```
### 6. Delete Folder
All the contents inside the specified folder will be deleted.
```bash
curl 'http://localhost:8080/deletefolder?folder=FolderName'
```
### 7.Folder Info
```bash
curl 'http://localhost:8080/folderinfo?folder=FolderName'
```
### 8. Zip Folder
The zipping process is performed in background.
```bash
curl 'http://localhost:8080/zip?folder=FolderName'
```
### 9. Zip Status
```bash
curl 'http://localhost:8080/zipstatus?folder=FolderName'
```
### 10. Download Zip
Download the zipped folder by specifying the folder which is zipped
```bash
curl 'http://localhost:8080/zipdownload?folder=FolderName'
```
