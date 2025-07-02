package logger

import (
	"log"
	"os"
)

type Logger interface {
	Info(message string)
	Error(message string)
	Debug(message string)
}

type SimpleLogger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
}

func New() *SimpleLogger {
	return &SimpleLogger{
		infoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		debugLogger: log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *SimpleLogger) Info(message string) {
	l.infoLogger.Println(message)
}

func (l *SimpleLogger) Error(message string) {
	l.errorLogger.Println(message)
}

func (l *SimpleLogger) Debug(message string) {
	l.debugLogger.Println(message)
}