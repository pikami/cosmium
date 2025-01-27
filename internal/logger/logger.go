package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
)

type LogLevelType int

var (
	LogLevelDebug  LogLevelType = 0
	LogLevelInfo   LogLevelType = 1
	LogLevelError  LogLevelType = 2
	LogLevelSilent LogLevelType = 10
)

type LogWriter struct {
	WriterLevel LogLevelType
}

var logLevelMutex sync.RWMutex
var logLevel = LogLevelInfo

var DebugLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
var InfoLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
var ErrorLogger = log.New(os.Stderr, "", log.Ldate|log.Ltime)

func DebugLn(v ...any) {
	if GetLogLevel() <= LogLevelDebug {
		prefix := getCallerPrefix()
		DebugLogger.Println(append([]interface{}{prefix}, v...)...)
	}
}

func Debug(v ...any) {
	if GetLogLevel() <= LogLevelDebug {
		prefix := getCallerPrefix()
		DebugLogger.Println(append([]interface{}{prefix}, v...)...)
	}
}

func Debugf(format string, v ...any) {
	if GetLogLevel() <= LogLevelDebug {
		prefix := getCallerPrefix()
		DebugLogger.Printf(prefix+format, v...)
	}
}

func InfoLn(v ...any) {
	if GetLogLevel() <= LogLevelInfo {
		InfoLogger.Println(v...)
	}
}

func Info(v ...any) {
	if GetLogLevel() <= LogLevelInfo {
		InfoLogger.Print(v...)
	}
}

func Infof(format string, v ...any) {
	if GetLogLevel() <= LogLevelInfo {
		InfoLogger.Printf(format, v...)
	}
}

func ErrorLn(v ...any) {
	if GetLogLevel() <= LogLevelError {
		prefix := getCallerPrefix()
		ErrorLogger.Println(append([]interface{}{prefix}, v...)...)
	}
}

func Error(v ...any) {
	if GetLogLevel() <= LogLevelError {
		prefix := getCallerPrefix()
		ErrorLogger.Print(append([]interface{}{prefix}, v...)...)
	}
}

func Errorf(format string, v ...any) {
	if GetLogLevel() <= LogLevelError {
		prefix := getCallerPrefix()
		ErrorLogger.Printf(prefix+format, v...)
	}
}

func (lw *LogWriter) Write(p []byte) (n int, err error) {
	switch lw.WriterLevel {
	case LogLevelDebug:
		Debug(string(p))
	case LogLevelInfo:
		Info(string(p))
	case LogLevelError:
		Error(string(p))
	}

	return len(p), nil
}

func ErrorWriter() *LogWriter {
	return &LogWriter{WriterLevel: LogLevelError}
}

func InfoWriter() *LogWriter {
	return &LogWriter{WriterLevel: LogLevelInfo}
}

func DebugWriter() *LogWriter {
	return &LogWriter{WriterLevel: LogLevelDebug}
}

func SetLogLevel(level LogLevelType) {
	logLevelMutex.Lock()
	defer logLevelMutex.Unlock()
	logLevel = level
}

func GetLogLevel() LogLevelType {
	logLevelMutex.RLock()
	defer logLevelMutex.RUnlock()
	return logLevel
}

func getCallerPrefix() string {
	_, file, line, ok := runtime.Caller(2)
	if ok {
		parts := strings.Split(file, "/")
		if len(parts) > 0 {
			file = parts[len(parts)-1]
		}
		return fmt.Sprintf("%s:%d - ", file, line)
	}

	return ""
}
