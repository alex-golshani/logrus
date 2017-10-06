package logrus

import (
	"bytes"
	"testing"
)

func BenchmarkLogEntryWithFieldsNoLog(b *testing.B) {
	logger := New(InfoLevel)
	logger.out = &bytes.Buffer{}
	entry := NewEntry(logger)
	for i := 0; i <= b.N; i++ {
		entry.AsDebug().WithField("test", "test").Write("Message")
	}
}

func BenchmarkLogEntryWithFieldsLogJSON(b *testing.B) {
	logger := New(DebugLevel)
	logger.SetOutput(&bytes.Buffer{})
	logger.SetFormatter(new(JSONFormatter))
	entry := NewEntry(logger)
	for i := 0; i <= b.N; i++ {
		entry.AsDebug().WithField("test", "test").Write("Message")
	}
}

func BenchmarkLogEntryWithFieldsLogText(b *testing.B) {
	logger := New(DebugLevel)
	logger.SetOutput(&bytes.Buffer{})
	entry := NewEntry(logger)
	for i := 0; i <= b.N; i++ {
		entry.AsDebug().WithField("test", "test").Write("Message")
	}
}
