package logrus

import (
	"fmt"
	"strings"
)

// Fields type, used to pass to `WithFields`.
type Fields map[string]interface{}

// level type
type Level uint32

// Convert the level to a string. E.g. PanicLevel becomes "panic".
func (level Level) String() string {
	switch level {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warning"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	case PanicLevel:
		return "panic"
	}

	return "unknown"
}

// ParseLevel takes a string level and returns the Logrus log level constant.
func ParseLevel(lvl string) (Level, error) {
	switch strings.ToLower(lvl) {
	case "panic":
		return PanicLevel, nil
	case "fatal":
		return FatalLevel, nil
	case "error":
		return ErrorLevel, nil
	case "warn", "warning":
		return WarnLevel, nil
	case "info":
		return InfoLevel, nil
	case "debug":
		return DebugLevel, nil
	}

	var l Level
	return l, fmt.Errorf("not a valid logrus level: %q", lvl)
}

// A constant exposing all logging levels
var AllLevels = []Level{
	PanicLevel,
	FatalLevel,
	ErrorLevel,
	WarnLevel,
	InfoLevel,
	DebugLevel,
}

// These are the different logging levels. You can set the logging level to log
// on your instance of Logger, obtained with `logrus.New()`.
const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// Message passed to Debug, Info, ...
	PanicLevel Level = iota
	// FatalLevel level. Logs and then calls `os.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
)

// The FieldLogger interface generalizes the Entry and Logger types
type FieldLogger interface {
	WithField(key string, value interface{}) *Entry
	WithFields(fields Fields) *Entry
	WithError(err error) *Entry

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	Debug(args ...interface{})
	Info(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})

	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Warningln(args ...interface{})
	Errorln(args ...interface{})
	Fatalln(args ...interface{})
	Panicln(args ...interface{})
}
