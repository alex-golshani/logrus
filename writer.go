package logrus

import (
	"bufio"
	"io"
	"runtime"
)

func (logger *Logger) Writer() *io.PipeWriter {
	return logger.WriterLevel(InfoLevel)
}

func (logger *Logger) WriterLevel(level Level) *io.PipeWriter {
	return NewEntry(logger).WriterLevel(level)
}

func (entry *Entry) Writer() *io.PipeWriter {
	return entry.WriterLevel(InfoLevel)
}

func (entry *Entry) WriterLevel(level Level) *io.PipeWriter {
	reader, writer := io.Pipe()

	var printFunc func(args ...interface{})

	switch level {
	case DebugLevel:
		printFunc = entry.AsDebug().Write
	case InfoLevel:
		printFunc = entry.AsInfo().Write
	case WarnLevel:
		printFunc = entry.AsWarning().Write
	case ErrorLevel:
		printFunc = entry.AsError().Write
	case FatalLevel:
		printFunc = entry.AsFatal().Write
	case PanicLevel:
		printFunc = entry.AsPanic().Write
	default:
		printFunc = entry.AsInfo().Write
	}

	go entry.writerScanner(reader, printFunc)
	runtime.SetFinalizer(writer, writerFinalizer)

	return writer
}

func (entry *Entry) writerScanner(reader *io.PipeReader, printFunc func(args ...interface{})) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		printFunc(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		entry.AsError().Write("Error while reading from Writer: %s", err)
	}
	reader.Close()
}

func writerFinalizer(writer *io.PipeWriter) {
	writer.Close()
}
