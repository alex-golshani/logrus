package logrus

import (
	"bytes"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"encoding/json"
)
func writeLogAndAssertJSON(loggerLevel Level, log func(*Logger), assertions func(Fields, *Logger)) {
	var buffer bytes.Buffer
	var fields Fields

	logger := New(loggerLevel)
	logger.Out = &buffer
	logger.formatter = new(JSONFormatter)

	log(logger)

	json.Unmarshal(buffer.Bytes(), &fields)

	assertions(fields, logger)
}

func writeLogAndAssertText(t *testing.T, loggerLevel Level, log func(*Logger), assertions func(Fields, *Logger)) {
	t.Helper()

	var buffer bytes.Buffer
	logger := New(loggerLevel)
	logger.Out = &buffer
	logger.formatter = &TextFormatter{
		DisableColors: true,
	}

	log(logger)

	fields := make(Fields)
	for _, kv := range strings.Split(strings.TrimSpace(buffer.String()), " ") {
		if !strings.Contains(kv, "=") {
			continue
		}
		kvArr := strings.Split(kv, "=")

		key := strings.TrimSpace(kvArr[0])
		val := kvArr[1]
		if kvArr[1][0] == '"' {
			var err error
			val, err = strconv.Unquote(val)
			assert.NoError(t, err)
		}
		fields[key] = val
	}

	assertions(fields, logger)
}

func TestLoggerDebugText(t *testing.T) {
	testCases := []struct {
		title       string
		loggerLevel Level
		message     string
		shouldLog   bool
	}{
		{
			title:       "logging_with_the_same_level_as_log_level_should_log",
			loggerLevel: DebugLevel,
			message:     "Message",
			shouldLog:   true,
		},
		{
			title:       "logging_with_the_level_higher_than_the_log_level_should_not_log",
			loggerLevel: InfoLevel,
			message:     "Message",
			shouldLog:   false,
		},
	}
	assrt := assert.New(t)
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			writeLogAndAssertText(t, tc.loggerLevel,
				func(lw *Logger) {
					lw.Debug(tc.message)
				},
				func(fields Fields, lw *Logger) {
					assrt.Equal(tc.loggerLevel, lw.level)
					msg, ok := fields[messageKey]
					if tc.shouldLog {
						if !ok {
							t.Error("Failed to retrieve the Message. Nothing was logged")
						}
						if logged, ok := checkLoggedField(tc.message, msg); !ok {
							t.Errorf("expected %s, received '%v'", tc.message, logged)
						}
						return
					}
					if ok {
						t.Errorf("we shouldn't have logged anything but the output was %v", fields)
					}
				})
		})
	}
}

func TestLoggerDebugJSON(t *testing.T) {
	testCases := []struct {
		title       string
		loggerLevel Level
		message     string
		shouldLog   bool
	}{
		{
			title:       "logging_with_the_same_level_as_log_level_should_log",
			loggerLevel: DebugLevel,
			message:     "Message",
			shouldLog:   true,
		},
		{
			title:       "logging_with_the_level_higher_than_the_log_level_should_not_log",
			loggerLevel: InfoLevel,
			message:     "Message",
			shouldLog:   false,
		},
	}
	assrt := assert.New(t)
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			writeLogAndAssertJSON(tc.loggerLevel,
				func(lw *Logger) {
					lw.Debug(tc.message)
				},
				func(fields Fields, lw *Logger) {
					assrt.Equal(tc.loggerLevel, lw.level)
					msg, ok := fields[messageKey]
					if tc.shouldLog {
						if !ok {
							t.Error("Failed to retrieve the Message. Nothing was logged")
						}
						if logged, ok := checkLoggedField(tc.message, msg); !ok {
							t.Errorf("expected %s, received '%v'", tc.message, logged)
						}
						return
					}
					if ok {
						t.Errorf("we shouldn't have logged anything but the output was %v", fields)
					}
				})
		})
	}
}

func TestLoggerInfoText(t *testing.T) {
	testCases := []struct {
		title       string
		loggerLevel Level
		message     string
		shouldLog   bool
	}{
		{
			title:       "logging_with_the_same_level_as_log_level_should_log",
			loggerLevel: DebugLevel,
			message:     "Message",
			shouldLog:   true,
		},
		{
			title:       "logging_with_the_level_higher_than_the_log_level_should_not_log",
			loggerLevel: WarnLevel,
			message:     "Message",
			shouldLog:   false,
		},
	}
	assrt := assert.New(t)
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			writeLogAndAssertText(t, tc.loggerLevel,
				func(lw *Logger) {
					lw.Info(tc.message)
				},
				func(fields Fields, lw *Logger) {
					assrt.Equal(tc.loggerLevel, lw.level)
					msg, ok := fields[messageKey]
					if tc.shouldLog {
						if !ok {
							t.Error("Failed to retrieve the Message. Nothing was logged")
						}
						if logged, ok := checkLoggedField(tc.message, msg); !ok {
							t.Errorf("expected %s, received '%v'", tc.message, logged)
						}
						return
					}
					if ok {
						t.Errorf("we shouldn't have logged anything but the output was %v", fields)
					}
				})
		})
	}
}

func TestLoggerInfoJSON(t *testing.T) {
	testCases := []struct {
		title       string
		loggerLevel Level
		message     string
		shouldLog   bool
	}{
		{
			title:       "logging_with_the_same_level_as_log_level_should_log",
			loggerLevel: DebugLevel,
			message:     "Message",
			shouldLog:   true,
		},
		{
			title:       "logging_with_the_level_higher_than_the_log_level_should_not_log",
			loggerLevel: WarnLevel,
			message:     "Message",
			shouldLog:   false,
		},
	}
	assrt := assert.New(t)
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			writeLogAndAssertJSON(tc.loggerLevel,
				func(lw *Logger) {
					lw.Info(tc.message)
				},
				func(fields Fields, lw *Logger) {
					assrt.Equal(tc.loggerLevel, lw.level)
					msg, ok := fields[messageKey]
					if tc.shouldLog {
						if !ok {
							t.Error("Failed to retrieve the Message. Nothing was logged")
						}
						if logged, ok := checkLoggedField(tc.message, msg); !ok {
							t.Errorf("expected %s, received '%v'", tc.message, logged)
						}
						return
					}
					if ok {
						t.Errorf("we shouldn't have logged anything but the output was %v", fields)
					}
				})
		})
	}
}
