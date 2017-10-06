package logrus

import (
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"
)

type Logger struct {
	// out The logs are `io.Copy`'d to this in a mutex. It's common to set this to a
	// file, or leave it default which is `os.Stderr`. You can also set this to
	// something more adventurous, such as logging to Kafka.
	out io.Writer
	// hooks for the Logger instance. These allow firing events based on logging
	// levels and log entries. For example, to send errors to an error tracking
	// service, log to StatsD or dump the core on fatal errors.
	hooks LevelHooks
	// formatter all log entries pass through the formatter before logged to out. The
	// included formatters are `TextFormatter` and `JSONFormatter` for which
	// TextFormatter is the default. In development (when a TTY is attached) it
	// logs with colors, but to a file it wouldn't. You can easily implement your
	// own that implements the `formatter` interface, see the `README` or included
	// formatters for examples.
	formatter Formatter
	// level the logging level the Logger should log at. This is typically (and defaults
	// to) `logrus.Info`, which allows Info(), Warn(), Error() and Fatal() to be
	// logged.
	level Level

	// MutexWrap used to sync writing to the log. Locking is enabled by Default
	mux MutexWrap

	// Reusable empty entry
	entryPool sync.Pool
}

// New creates a new instance of Logger. Configuration should be set by changing `formatter` (default TextFormatter),
// `out` (default os.Stderr) and `hooks` directly on the default Logger instance.
// It's recommended to make this a global instance called `log`.
func New(level Level) *Logger {
	return &Logger{
		out: os.Stderr,
		formatter: &TextFormatter{
			DisableSorting: true,
		},
		hooks: make(LevelHooks),
		level: level,
	}
}

// AsLevel creates a new entry and sets the level to the specified value.
// Make sure you call this method before calling WithField, WithFields and WithError methods
func (logger *Logger) AsLevel(level Level) *Entry {
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.AsLevel(level)
}

// AsDebug creates a new entry and sets the level to `debug`
// Make sure you call this method before calling WithField, WithFields and WithError methods
func (logger *Logger) AsDebug() *Entry {
	return logger.AsLevel(DebugLevel)
}

// AsInfo creates a new entry and sets the level to `info`
// Make sure you call this method before calling WithField, WithFields and WithError methods
func (logger *Logger) AsInfo() *Entry {
	return logger.AsLevel(InfoLevel)
}

// AsWarning creates a new entry and sets the level to `warning`
// Make sure you call this method before calling WithField, WithFields and WithError methods
func (logger *Logger) AsWarning() *Entry {
	return logger.AsLevel(WarnLevel)
}

// AsError creates a new entry and sets the level to `error`
// Make sure you call this method before calling WithField, WithFields and WithError methods
func (logger *Logger) AsError() *Entry {
	return logger.AsLevel(ErrorLevel)
}

// AsFatal creates a new entry and sets the level to `fatal`
// Make sure you call this method before calling WithField, WithFields and WithError methods
func (logger *Logger) AsFatal() *Entry {
	return logger.AsLevel(FatalLevel)
}

// AsPanic creates a new entry and sets the level to `panic`
// Make sure you call this method before calling WithField, WithFields and WithError methods
func (logger *Logger) AsPanic() *Entry {
	return logger.AsLevel(PanicLevel)
}

// WithFields creates a new log entry object and adds a struct of fields to the entry
func (logger *Logger) WithFields(fields Fields) *Entry {
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WithFields(fields)
}

// WithField creates a new log entry object and adds a field to the entry.
//If you want multiple fields, use `WithFields`.
func (logger *Logger) WithField(key string, value interface{}) *Entry {
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WithField(key, value)
}

// WithError adds an error as single field to the log entry
func (logger *Logger) WithError(err error) *Entry {
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WithError(err)
}

// Debugf logs a formatted string at debug level
func (logger *Logger) Debugf(format string, args ...interface{}) {
	logger.log(DebugLevel, formatted, format, args...)
}

// Infof logs a formatted string at info level
func (logger *Logger) Infof(format string, args ...interface{}) {
	logger.log(InfoLevel, formatted, format, args...)
}

// Warningf logs a formatted string at warning level
func (logger *Logger) Warningf(format string, args ...interface{}) {
	logger.log(WarnLevel, formatted, format, args...)
}

// Errorf logs a formatted string at error level
func (logger *Logger) Errorf(format string, args ...interface{}) {
	logger.log(ErrorLevel, formatted, format, args...)
}

// Fatalf logs a formatted string at fatal level and terminated the execution of the application with exit code 1
func (logger *Logger) Fatalf(format string, args ...interface{}) {
	logger.log(FatalLevel, formatted, format, args...)
	Exit(1)
}

// Panicf logs a formatted string at panic level and terminates the execution of the application with a panic
func (logger *Logger) Panicf(format string, args ...interface{}) {
	logger.log(PanicLevel, formatted, format, args...)
	panic(fmt.Sprint(args...))
}

// Debug logs a message at debug level
func (logger *Logger) Debug(args ...interface{}) {
	logger.log(DebugLevel, unformatted, "", args...)
}

// Info logs a message at info level
func (logger *Logger) Info(args ...interface{}) {
	logger.log(InfoLevel, unformatted, "", args...)
}

// Warning logs a message at warning level
func (logger *Logger) Warning(args ...interface{}) {
	logger.log(WarnLevel, unformatted, "", args...)
}

// Error logs a message at error level
func (logger *Logger) Error(args ...interface{}) {
	logger.log(ErrorLevel, unformatted, "", args...)
}

// Fatal logs a message at fatal level and terminates the app with exit code 1
func (logger *Logger) Fatal(args ...interface{}) {
	logger.log(FatalLevel, unformatted, "", args...)
	Exit(1)
}

// Panic logs a message at panic level and terminates the execution of the application with a panic
func (logger *Logger) Panic(args ...interface{}) {
	msg := fmt.Sprint(args...)
	logger.log(PanicLevel, unformatted, "", args...)
	panic(msg)
}

// Debugln logs a message followed by a new line at debug level
func (logger *Logger) Debugln(args ...interface{}) {
	logger.log(DebugLevel, newLine, "", args...)
}

// Infoln logs a message followed by a new line at info level
func (logger *Logger) Infoln(args ...interface{}) {
	logger.log(InfoLevel, newLine, "", args...)
}

// Warningln logs a message followed by a new line at warning level
func (logger *Logger) Warningln(args ...interface{}) {
	logger.log(WarnLevel, newLine, "", args...)
}

// Errorln logs a message followed by a new line at error level
func (logger *Logger) Errorln(args ...interface{}) {
	logger.log(ErrorLevel, newLine, "", args...)
}

//Fatalln logs a message followed by a new line at fatal level and terminates the execution of the application with exit code 1
func (logger *Logger) Fatalln(args ...interface{}) {
	logger.log(FatalLevel, newLine, "", args...)
	Exit(1)
}

// Panicln logs a message followed by a new line at panic level and terminates the execution of the application with a panic
func (logger *Logger) Panicln(args ...interface{}) {
	msg := fmt.Sprint(args...)
	logger.log(PanicLevel, newLine, "", msg)
	panic(msg)
}

//SetNoLock when file is opened with appending mode, it's safe to
//write concurrently to a file (within 4k Message on Linux).
//In this case user can choose to disable the lock.
func (logger *Logger) SetNoLock() {
	logger.mux.Disable()
}

// SetOutput sets the Logger's output.
func (logger *Logger) SetOutput(out io.Writer) {
	logger.mux.Lock()
	defer logger.mux.Unlock()
	logger.out = out
}

// SetFormatter sets the Logger's formatter.
func (logger *Logger) SetFormatter(formatter Formatter) {
	logger.mux.Lock()
	defer logger.mux.Unlock()
	logger.formatter = formatter
}

// AddHook adds an external hook to the logger
// The hooks will get executed before we log the entry
func (logger *Logger) AddHook(hook Hook) {
	logger.mux.Lock()
	defer logger.mux.Unlock()
	logger.hooks.Add(hook)
}

func (logger *Logger) releaseEntry(entry *Entry) {
	logger.entryPool.Put(entry)
}

func (logger *Logger) log(level Level, mode formatMode, format string, args ...interface{}) {
	if logger.getLevel() >= level {
		entry := logger.newEntry()
		message := constructMessage(mode, format, args...)
		entry.Level = level
		entry.log(message)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) newEntry() *Entry {
	entry, ok := logger.entryPool.Get().(*Entry)
	if ok {
		return entry
	}
	return NewEntry(logger)
}

func (logger *Logger) getLevel() Level {
	return Level(atomic.LoadUint32((*uint32)(&logger.level)))
}
