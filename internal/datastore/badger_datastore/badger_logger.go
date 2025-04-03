package badgerdatastore

import (
	"github.com/dgraph-io/badger/v4"
	"github.com/pikami/cosmium/internal/logger"
)

type badgerLogger struct{}

func newBadgerLogger() badger.Logger {
	return &badgerLogger{}
}

func (l *badgerLogger) Errorf(format string, v ...interface{}) {
	logger.Errorf(format, v...)
}

func (l *badgerLogger) Warningf(format string, v ...interface{}) {
	logger.Infof(format, v...)
}

func (l *badgerLogger) Infof(format string, v ...interface{}) {
	logger.Infof(format, v...)
}

func (l *badgerLogger) Debugf(format string, v ...interface{}) {
	logger.Debugf(format, v...)
}
