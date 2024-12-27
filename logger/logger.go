package logger

import (
	"fileServer/constants"
	"fmt"
	"runtime"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type BufferedLogger struct {
	LogChannel  chan LogEntry
	FlushTicker *time.Ticker
	Done        chan bool
	Wg          sync.WaitGroup
	Logger      *log.Logger
}

type LogEntry struct {
	LogLevel log.Level
	LogField log.Fields
	Message  string
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
		LogChannel:  make(chan LogEntry, constants.LogBufferSize),
		FlushTicker: time.NewTicker(flushInterval),
		Done:        make(chan bool),
		Logger:      logger,
	}
	bufferedLogger.Wg.Add(1)
	go bufferedLogger.startFlushTicker()
	return bufferedLogger
}

func (bl *BufferedLogger) Log(level log.Level, field log.Fields, message string) {
	bl.LogChannel <- LogEntry{LogLevel: level, LogField: field, Message: message}
}

func (bl *BufferedLogger) startFlushTicker() {
	defer bl.Wg.Done()
	for {
		select {
		case <-bl.FlushTicker.C:
			bl.FlushLogs()
		case <-bl.Done:
			bl.FlushLogs()
			return
		}
	}
}

func (bl *BufferedLogger) FlushLogs() {
	for {
		select {
		case LogEntry := <-bl.LogChannel:
			switch LogEntry.LogLevel {
			case log.TraceLevel:
				bl.Logger.WithFields(LogEntry.LogField).Trace(LogEntry.Message)
			case log.DebugLevel:
				bl.Logger.WithFields(LogEntry.LogField).Debug(LogEntry.Message)
			case log.InfoLevel:
				bl.Logger.WithFields(LogEntry.LogField).Info(LogEntry.Message)
			case log.WarnLevel:
				bl.Logger.WithFields(LogEntry.LogField).Warn(LogEntry.Message)
			case log.ErrorLevel:
				bl.Logger.WithFields(LogEntry.LogField).Error(LogEntry.Message)
			case log.FatalLevel:
				bl.Logger.WithFields(LogEntry.LogField).Fatal(LogEntry.Message)
			case log.PanicLevel:
				bl.Logger.WithFields(LogEntry.LogField).Panic(LogEntry.Message)
			default:
				bl.Logger.WithFields(LogEntry.LogField).Info(LogEntry.Message)
			}
		default:
			return
		}
	}
}