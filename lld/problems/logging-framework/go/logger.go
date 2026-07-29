// Package loggingframework implements a small, pluggable logging framework.
//
// A Logger holds a minimum severity threshold and a list of Appenders
// (Strategy pattern) that it dispatches formatted records to. Appenders are
// interchangeable output destinations (console, file, ...); new ones can be
// added without touching Logger itself.
package loggingframework

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Level is a log severity, ordered DEBUG < INFO < WARN < ERROR.
type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Record is a single formatted log entry handed to every Appender.
type Record struct {
	Timestamp time.Time
	Level     Level
	Message   string
}

// Format renders the record as "<RFC3339> [LEVEL] message".
func (r Record) Format() string {
	return fmt.Sprintf("%s [%s] %s", r.Timestamp.Format(time.RFC3339), r.Level, r.Message)
}

// Appender is a pluggable output destination for log records.
type Appender interface {
	Append(record Record)
}

// ConsoleAppender writes formatted records to an io.Writer (os.Stdout by
// default).
type ConsoleAppender struct {
	out io.Writer
}

// NewConsoleAppender returns a ConsoleAppender writing to os.Stdout.
func NewConsoleAppender() *ConsoleAppender {
	return &ConsoleAppender{out: os.Stdout}
}

// NewConsoleAppenderWithWriter returns a ConsoleAppender writing to w, useful
// for tests that want to capture output without touching stdout.
func NewConsoleAppenderWithWriter(w io.Writer) *ConsoleAppender {
	return &ConsoleAppender{out: w}
}

func (c *ConsoleAppender) Append(record Record) {
	fmt.Fprintln(c.out, record.Format())
}

// FileAppender appends formatted records, one per line, to a file at a
// caller-supplied path. The file is opened once and kept open for the
// lifetime of the appender.
type FileAppender struct {
	mu   sync.Mutex
	file *os.File
}

// NewFileAppender opens (creating/truncating-safe append mode) the file at
// path for appending and returns a FileAppender backed by it.
func NewFileAppender(path string) (*FileAppender, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("open file appender: %w", err)
	}
	return &FileAppender{file: f}, nil
}

func (a *FileAppender) Append(record Record) {
	a.mu.Lock()
	defer a.mu.Unlock()
	fmt.Fprintln(a.file, record.Format())
}

// Close closes the underlying file.
func (a *FileAppender) Close() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.file.Close()
}

// Logger dispatches records to its appenders when the record's level meets
// or exceeds the logger's minimum threshold. Safe for concurrent use.
type Logger struct {
	mu        sync.Mutex
	level     Level
	appenders []Appender
}

// NewLogger returns a Logger with the given minimum level and appenders.
func NewLogger(level Level, appenders ...Appender) *Logger {
	return &Logger{
		level:     level,
		appenders: append([]Appender{}, appenders...),
	}
}

// AddAppender registers an additional appender.
func (l *Logger) AddAppender(a Appender) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.appenders = append(l.appenders, a)
}

// SetLevel updates the logger's minimum severity threshold.
func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// Log formats a record and dispatches it to every appender, provided level
// is at or above the logger's threshold.
func (l *Logger) Log(level Level, message string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if level < l.level {
		return
	}
	record := Record{Timestamp: time.Now(), Level: level, Message: message}
	for _, a := range l.appenders {
		a.Append(record)
	}
}

func (l *Logger) Debug(message string) { l.Log(DEBUG, message) }
func (l *Logger) Info(message string)  { l.Log(INFO, message) }
func (l *Logger) Warn(message string)  { l.Log(WARN, message) }
func (l *Logger) Error(message string) { l.Log(ERROR, message) }
