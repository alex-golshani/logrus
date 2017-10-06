package logrus

import (
	"fmt"
	"os"
	"time"
)

// Entry an entry is the final or intermediate Logrus logging entry. It contains all
// the fields passed with WithField{,s}. It's finally logged when Write method is called. These objects can be reused and
// passed around as much as you wish to avoid field duplication.
type Entry struct {
	Logger *Logger

	// Data contains all the fields set by the user.
	Data Fields

	// Time at which the log entry was created
	Time time.Time

	// Level the log entry was logged at: Debug, Info, Warn, Error, Fatal or Panic
	Level Level

	// Message passed to Write method
	Message string
}

// NewEntry creates a new log entry
func NewEntry(logger *Logger) *Entry {
	// Default is three fields, give a little extra room
	return newLogEntry(logger, logger.level, make(Fields, 5))
}

// NewEntryWithFields creates a new log entry and adds a struct of fields to the entry
func NewEntryWithFields(logger *Logger, fields Fields) *Entry {
	return newLogEntry(logger, logger.level, fields)
}

// NewEntryWithField creates a new log entry and adds a field to the entry
//If you want multiple fields, use `NewEntryWithFields`
func NewEntryWithField(logger *Logger, key string, value interface{}) *Entry {
	//Do not change this to Fields{key:value}. You will end up getting more allocations
	fields := make(Fields, 1)
	fields[key] = value
	return newLogEntry(logger, logger.level, fields)
}

// AsLevel clones the entry into a new log entry and sets the level to the specified value.
// Make sure you call this method before calling WithField, WithFields and WithError methods
func (entry *Entry) AsLevel(level Level) *Entry {
	return newLogEntry(entry.Logger, level, entry.Data)
}

// AsDebug clones the entry into a new log entry and sets the level to `debug`
// Make sure you call this method before calling WithField, WithFields and WithError methods
func (entry *Entry) AsDebug() *Entry {
	return entry.AsLevel(DebugLevel)
}

// AsInfo clones the entry into a new log entry and sets the level to `info`
// Make sure you call this method before calling WithField, WithFields and WithError methods
func (entry *Entry) AsInfo() *Entry {
	return entry.AsLevel(InfoLevel)
}

// AsWarning clones the entry into a new log entry and sets the level to `warning`
// Make sure you call this method before calling WithField, WithFields and WithError methods
func (entry *Entry) AsWarning() *Entry {
	return entry.AsLevel(WarnLevel)
}

// AsError clones the entry into a new log entry and sets the level to `error`
// Make sure you call this method before calling WithField, WithFields and WithError methods
func (entry *Entry) AsError() *Entry {
	return entry.AsLevel(ErrorLevel)
}

// AsFatal clones the entry into a new log entry and sets the level to `fatal`
// Make sure you call this method before calling WithField, WithFields and WithError methods
func (entry *Entry) AsFatal() *Entry {
	return entry.AsLevel(FatalLevel)
}

// AsPanic clones the entry into a new log entry and sets the level to `panic`
// Make sure you call this method before calling WithField, WithFields and WithError methods
func (entry *Entry) AsPanic() *Entry {
	return entry.AsLevel(PanicLevel)
}

// WithField adds a field to the log entry, note that it doesn't log until you call Write.
func (entry *Entry) WithField(key string, value interface{}) *Entry {
	if entry.Level > entry.Logger.getLevel() {
		return entry
	}
	//Do not change this to Fields{key:value}. You will end up getting more allocations
	fields := make(Fields, 1)
	fields[key] = value
	return entry.WithFields(fields)
}

// WithFields adds a struct of fields to the log entry
func (entry *Entry) WithFields(fields Fields) *Entry {
	if entry.Level > entry.Logger.getLevel() {
		return entry
	}
	data := make(Fields, len(entry.Data)+len(fields))
	for k, v := range entry.Data {
		data[k] = v
	}
	for k, v := range fields {
		data[k] = v
	}
	return newLogEntry(entry.Logger, entry.Level, data)
}

// WithError adds an error as single field to the log entry
func (entry *Entry) WithError(err error) *Entry {
	return entry.WithField(errorKey, err)
}

func (entry *Entry) Writef(format string, args ...interface{}) {
	entry.write(formatted, format, args...)
}

func (entry *Entry) Write(args ...interface{}) {
	entry.write(unformatted, "", args...)
}

func (entry *Entry) Writeln(args ...interface{}) {
	entry.write(newLine, "", args...)
}

func newLogEntry(logger *Logger, level Level, data Fields) *Entry {
	return &Entry{
		Logger: logger,
		Data:   data,
		Level:  level,
	}
}

func (entry *Entry) write(mode formatMode, format string, args ...interface{}) {
	if entry.Logger.getLevel() >= entry.Level {
		message := constructMessage(mode, format, args...)
		entry.log(message)
	}
}

func constructMessage(mode formatMode, format string, args ...interface{}) string {
	switch mode {
	case formatted:
		return fmt.Sprintf(format, args...)
	case unformatted:
		return fmt.Sprint(args...)
	case newLine:
		return sprintlnn(args...)
	}
	return fmt.Sprintf(format, args...)
}

func (entry *Entry) log(msg string) {
	entry.Time = time.Now()
	entry.Message = msg

	entry.Logger.mux.Lock()
	err := entry.Logger.hooks.Fire(entry.Level, entry)
	entry.Logger.mux.Unlock()
	if err != nil {
		entry.Logger.mux.Lock()
		fmt.Fprintf(os.Stderr, "Failed to fire the hook: %v\n", err)
		entry.Logger.mux.Unlock()
	}
	serialized, err := entry.Logger.formatter.Format(entry)
	if err != nil {
		entry.Logger.mux.Lock()
		fmt.Fprintf(os.Stderr, "Failed to obtain reader, %v\n", err)
		entry.Logger.mux.Unlock()
	} else {
		entry.Logger.mux.Lock()
		_, err = entry.Logger.out.Write(serialized)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write to log, %v\n", err)
		}
		entry.Logger.mux.Unlock()
	}

	if entry.Level == FatalLevel {
		Exit(1)
	}

	// To avoid Entry#log() returning a value that only would make sense for
	// panic() to use in Entry#Panic(), we avoid the allocation by checking
	// directly here.
	if entry.Level <= PanicLevel {
		panic(&entry)
	}
}

// String returns the string representation from the reader and ultimately the
// formatter.
func (entry *Entry) String() (string, error) {
	serialized, err := entry.Logger.formatter.Format(entry)
	if err != nil {
		return "", err
	}
	str := string(serialized)
	return str, nil
}

// sprintlnn => Sprint no newline. This is to get the behavior of how
// fmt.Sprintln where spaces are always added between operands, regardless of
// their type. Instead of vendoring the Sprintln implementation to spare a
// string allocation, we do the simplest thing.
func sprintlnn(args ...interface{}) string {
	msg := fmt.Sprintln(args...)
	return msg[:len(msg)-1]
}
