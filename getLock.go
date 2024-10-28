package main

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

func getFileLock(filePath string) *sync.Mutex {
	lock, _ := fileLocks.LoadOrStore(filePath, &sync.Mutex{})
	logger.Log(log.InfoLevel, log.Fields{"method": "getFileLock"}, "File/Folder lock acquired")
	return lock.(*sync.Mutex)
}
