package logger

import (
	"log"
	"os"
)

type Logger struct {
	Error *log.Logger
	Debug *log.Logger
}

var Log = NewLogger()

func NewLogger() *Logger {
	return &Logger{
		Error: log.New(os.Stdout, "[error]: ", log.LstdFlags),
		Debug: log.New(os.Stdout, "[debug]: ", log.LstdFlags),
	}
}
