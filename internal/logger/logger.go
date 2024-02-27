package logger

import (
	"log"
	"os"

	"github.com/pikami/cosmium/api/config"
)

var DebugLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
var InfoLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
var ErrorLogger = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile)

func Debug(v ...any) {
	if config.Config.Debug {
		DebugLogger.Println(v...)
	}
}

func Debugf(format string, v ...any) {
	if config.Config.Debug {
		DebugLogger.Printf(format, v...)
	}
}

func Info(v ...any) {
	InfoLogger.Println(v...)
}

func Infof(format string, v ...any) {
	InfoLogger.Printf(format, v...)
}

func Error(v ...any) {
	ErrorLogger.Println(v...)
}

func Errorf(format string, v ...any) {
	ErrorLogger.Printf(format, v...)
}
