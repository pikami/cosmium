package logger

import (
	"log"
	"os"
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

var LogLevel = LogLevelInfo

var DebugLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
var InfoLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
var ErrorLogger = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile)

func DebugLn(v ...any) {
	if LogLevel <= LogLevelDebug {
		DebugLogger.Println(v...)
	}
}

func Debug(v ...any) {
	if LogLevel <= LogLevelDebug {
		DebugLogger.Print(v...)
	}
}

func Debugf(format string, v ...any) {
	if LogLevel <= LogLevelDebug {
		DebugLogger.Printf(format, v...)
	}
}

func InfoLn(v ...any) {
	if LogLevel <= LogLevelInfo {
		InfoLogger.Println(v...)
	}
}

func Info(v ...any) {
	if LogLevel <= LogLevelInfo {
		InfoLogger.Print(v...)
	}
}

func Infof(format string, v ...any) {
	if LogLevel <= LogLevelInfo {
		InfoLogger.Printf(format, v...)
	}
}

func ErrorLn(v ...any) {
	if LogLevel <= LogLevelError {
		ErrorLogger.Println(v...)
	}
}

func Error(v ...any) {
	if LogLevel <= LogLevelError {
		ErrorLogger.Print(v...)
	}
}

func Errorf(format string, v ...any) {
	if LogLevel <= LogLevelError {
		ErrorLogger.Printf(format, v...)
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
