package main

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	uploadDir       = "./uploads"
	logDir          = "./logs"
	logFile         = "logFile.log"
	workerCount     = 100
	loadWorkerCount = 100
	logBufferSize   = 200
)

var (
	fileLocks         sync.Map
	metaDataMap       sync.Map
	metadataLock      sync.Mutex
	folderMetadataMap sync.Map
	taskQueue         = make(chan string, 150)
	loadQueue         = make(chan string, 200)
	wg                sync.WaitGroup
	loadwg            sync.WaitGroup
	zipStatuses       sync.Map
	logger            *BufferedLogger
)

type Response struct {
	StatusCode   int       `json:"status_code"`
	Status       string    `json:"status"`
	Message      string    `json:"message"`
	ResponseTime time.Time `json:"response_time"`
}

type FileMetadata struct {
	FileName     string      `json:"file_name"`
	FilePath     string      `json:"file_path"`
	FolderPath   string      `json:"folder_path"`
	FileSize     int64       `json:"file_size"`
	ModifiedTime string      `json:"modified_time"`
	CreatedTime  string      `json:"created_time"`
	FileMode     fs.FileMode `json:"file_mode"`
	IsDirectory  bool        `json:"is_directory"`
}

type FolderMetadata struct {
	FolderName   string      `json:"folder_name"`
	FolderPath   string      `json:"folder_path"`
	FolderSize   int64       `json:"folder_size"`
	FilesCount   int         `json:"files_count"`
	ModifiedTime string      `json:"modified_time"`
	CreatedTime  string      `json:"created_time"`
	FolderMode   os.FileMode `json:"fileMode"`
	IsDirectory  bool        `json:"is_directory"`
}

type zipStatus struct {
	Status    string
	StartTime time.Time
	EndTime   time.Time
	FilePath  string
	ErrorMsg  string
}

type BufferedLogger struct {
	logChannel  chan LogEntry
	flushTicker *time.Ticker
	done        chan bool
	wg          sync.WaitGroup
	logger      *log.Logger
}

type LogEntry struct {
	LogLevel log.Level
	LogField log.Fields
	Message  string
}

func init() {
	_, err := os.Stat(uploadDir)
	if os.IsNotExist(err) {
		os.Mkdir(uploadDir, os.ModePerm)
	}

	logger = NewBufferedLogger(filepath.Join(logDir, logFile), 1, 1, time.Millisecond)
}

func main() {
	logField := log.Fields{
		"method": "main",
	}
	logger.Log(log.InfoLevel, logField, "Application starts...")
	loadWorkerPool()
	loadFileMetadataAtStart()
	loadFolderMetadata()
	workerPool()
	server := http.Server{
		Addr: ":8080",
	}
	http.HandleFunc("/upload", uploadFileHandler)
	http.HandleFunc("/download", downloadFileHandler)
	http.HandleFunc("/fileinfo", fileInfoHandler)
	http.HandleFunc("/delete", deleteFileHandler)
	http.HandleFunc("/folderinfo", folderInfoHandler)
	http.HandleFunc("/createfolder", folderCreateHandler)
	http.HandleFunc("/deletefolder", deleteFolderHandler)
	http.HandleFunc("/zip", zipFolderHandler)
	http.HandleFunc("/zipdownload", downloadZipHandler)
	http.HandleFunc("/zipstatus", zipStatusHandler)
	logger.Log(log.InfoLevel, logField, "Server is listening in http://localhost:8080")
	log.Println("Server is listening in http://localhost:8080")
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Log(log.FatalLevel, logField, fmt.Sprintf("HTTP shutdown error: %v", err))
		}
		logger.Log(log.InfoLevel, logField, "Server is shutting down")
	}()
	err := server.ListenAndServe()
	if err == http.ErrServerClosed {
		fmt.Println("Server is shutting down")
		Shutdown()
	}
}

func Shutdown() {
	close(taskQueue)
	wg.Wait()
	logger.Log(log.InfoLevel, log.Fields{"method": "Shutdown"}, "All task completed, Shutting down...")
	log.Println("All task completed, shutting down...")
}

func NewBufferedLogger(logFilePath string, maxSize int, maxAge int, flushInterval time.Duration) *BufferedLogger {
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    maxSize,
		MaxAge:     maxAge,
		MaxBackups: 3,
		Compress:   true,
	}
	logger := log.New()
	logger.SetOutput(lumberjackLogger)
	logger.SetFormatter(&log.JSONFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			return "", fmt.Sprintf("%s:%d", f.File, f.Line)
		},
	})
	logger.SetLevel(log.TraceLevel)
	bufferedLogger := &BufferedLogger{
		logChannel:  make(chan LogEntry, logBufferSize),
		flushTicker: time.NewTicker(flushInterval),
		done:        make(chan bool),
		logger:      logger,
	}
	bufferedLogger.wg.Add(1)
	go bufferedLogger.startFlushTicker()
	return bufferedLogger
}

func (bl *BufferedLogger) Log(level log.Level, field log.Fields, message string) {
	bl.logChannel <- LogEntry{LogLevel: level, LogField: field, Message: message}
}

func (bl *BufferedLogger) startFlushTicker() {
	defer bl.wg.Done()
	for {
		select {
		case <-bl.flushTicker.C:
			bl.FlushLogs()
		case <-bl.done:
			bl.FlushLogs()
			return
		}
	}
}

func (bl *BufferedLogger) FlushLogs() {
	for {
		select {
		case LogEntry := <-bl.logChannel:
			switch LogEntry.LogLevel {
			case log.TraceLevel:
				bl.logger.WithFields(LogEntry.LogField).Trace(LogEntry.Message)
			case log.DebugLevel:
				bl.logger.WithFields(LogEntry.LogField).Debug(LogEntry.Message)
			case log.InfoLevel:
				bl.logger.WithFields(LogEntry.LogField).Info(LogEntry.Message)
			case log.WarnLevel:
				bl.logger.WithFields(LogEntry.LogField).Warn(LogEntry.Message)
			case log.ErrorLevel:
				bl.logger.WithFields(LogEntry.LogField).Error(LogEntry.Message)
			case log.FatalLevel:
				bl.logger.WithFields(LogEntry.LogField).Fatal(LogEntry.Message)
			case log.PanicLevel:
				bl.logger.WithFields(LogEntry.LogField).Panic(LogEntry.Message)
			default:
				bl.logger.WithFields(LogEntry.LogField).Info(LogEntry.Message)
			}
		default:
			return
		}
	}
}
