package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

var (
	logger *log.Logger
	level  LogLevel
)

func init() {
	logger = log.New(os.Stdout, "", 0)
	level = INFO // default log level
}

func SetLogLevel(l LogLevel) {
	level = l
}

func logMessage(messageLevel LogLevel, format string, v ...interface{}) {
	if messageLevel >= level {
		prefix := time.Now().Format("2006-01-02 15:04:05")
		levelStr := ""
		switch messageLevel {
		case DEBUG:
			levelStr = "DEBUG"
		case INFO:
			levelStr = "INFO "
		case WARN:
			levelStr = "WARN "
		case ERROR:
			levelStr = "ERROR"
		}
		message := fmt.Sprintf(format, v...)
		logger.Printf("%s [%s] %s\n", prefix, levelStr, message)
	}
}

func Debug(format string, v ...interface{}) {
	logMessage(DEBUG, format, v...)
}

func Info(format string, v ...interface{}) {
	logMessage(INFO, format, v...)
}

func Warn(format string, v ...interface{}) {
	logMessage(WARN, format, v...)
}

func Error(format string, v ...interface{}) {
	logMessage(ERROR, format, v...)
}
