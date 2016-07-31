// Logger is a simple wrapper for logrus.Logger which extends it's capabilities by providing Logger's
// per category.  This is useful when an application is broken up into different components
// and different levels set per category
package logger

import (
	"github.com/Sirupsen/logrus"
	"os"
)

const (
	DefaultCategory = "_"
	DefaultTimeFormat = "2006-01-02T15:04:05.000"
)

// Logger is a wrapper for a logrus.Logger which adds a category
type Logger struct {
	// Current category for this logger
	category string
	// inherit from logrus logger
	logrus.Logger
}

var loggers map[string]*Logger
var globalLogger *Logger

func init() {
	loggers = map[string]*Logger{}
	globalLogger = GetLogger(DefaultCategory)
}

// GetLogger returns an existing logger for the specified category or
// creates a new one if it hasn't been defined
func GetLogger(category string) *Logger {
	if l, exists := loggers[category]; exists {
		return l
	}
	l := newLogger(category)
	loggers[category] = l
	return l
}

// Logger is a simple accessor to the global logging instance
func Logger() *Logger {
	return globalLogger
}

// SetLevel enforces the specified level on the specified category.  If a logger
// for the category hasn't been defined then this is a no-op.  In some cases
// you may want to globally initialize all logging modules ahead of time, in this
// scenario you can call GetLogger and set the returned pointers Level
func SetLevel(level logrus.Level, category string) {
	if l, exists := loggers[category]; exists {
		l.Level = level
	}
}

func newLogger(module string) *Logger {
	log := &Logger{module}
	log.Out = os.Stderr
	log.Formatter = logrus.TextFormatter{ DisableTimestamp: false, TimestampFormat: DefaultTimeFormat}
	log.Hooks = logrus.InfoLevel

	return log
}
