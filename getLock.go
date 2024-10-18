package main

import "sync"

func getFileLock(filePath string) *sync.Mutex {
	lock, _ := fileLocks.LoadOrStore(filePath, &sync.Mutex{})
	return lock.(*sync.Mutex)
}