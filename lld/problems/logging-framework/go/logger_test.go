package loggingframework

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

// mockAppender records every Record it receives; safe for concurrent use so
// it can be shared across goroutines in the thread-safety test.
type mockAppender struct {
	mu      sync.Mutex
	records []Record
}

func (m *mockAppender) Append(record Record) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.records = append(m.records, record)
}

func (m *mockAppender) count() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.records)
}

func TestLevelFiltering(t *testing.T) {
	tests := []struct {
		name      string
		threshold Level
		logAt     Level
		wantCount int
	}{
		{"debug threshold logs debug", DEBUG, DEBUG, 1},
		{"info threshold drops debug", INFO, DEBUG, 0},
		{"info threshold logs info", INFO, INFO, 1},
		{"warn threshold drops info", WARN, INFO, 0},
		{"warn threshold logs warn", WARN, WARN, 1},
		{"error threshold drops warn", ERROR, WARN, 0},
		{"error threshold logs error", ERROR, ERROR, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appender := &mockAppender{}
			logger := NewLogger(tt.threshold, appender)
			logger.Log(tt.logAt, "message")
			if got := appender.count(); got != tt.wantCount {
				t.Fatalf("expected %d records, got %d", tt.wantCount, got)
			}
		})
	}
}

func TestMultipleAppendersAllReceiveRecord(t *testing.T) {
	a1 := &mockAppender{}
	a2 := &mockAppender{}
	a3 := &mockAppender{}
	logger := NewLogger(INFO, a1, a2, a3)

	logger.Info("hello")
	logger.Debug("should be filtered")

	for i, a := range []*mockAppender{a1, a2, a3} {
		if got := a.count(); got != 1 {
			t.Fatalf("appender %d: expected 1 record, got %d", i, got)
		}
	}
}

func TestAddAppenderAfterConstruction(t *testing.T) {
	a1 := &mockAppender{}
	logger := NewLogger(DEBUG, a1)

	a2 := &mockAppender{}
	logger.AddAppender(a2)

	logger.Info("hi")
	if a1.count() != 1 || a2.count() != 1 {
		t.Fatalf("expected both appenders to receive the record, got a1=%d a2=%d", a1.count(), a2.count())
	}
}

func TestSetLevelChangesThreshold(t *testing.T) {
	appender := &mockAppender{}
	logger := NewLogger(ERROR, appender)

	logger.Warn("filtered")
	if appender.count() != 0 {
		t.Fatalf("expected warn to be filtered at ERROR threshold")
	}

	logger.SetLevel(WARN)
	logger.Warn("passes now")
	if appender.count() != 1 {
		t.Fatalf("expected warn to pass after lowering threshold, got %d", appender.count())
	}
}

func TestRecordFormatContainsLevelAndMessage(t *testing.T) {
	var buf bytes.Buffer
	appender := NewConsoleAppenderWithWriter(&buf)
	logger := NewLogger(DEBUG, appender)

	logger.Error("disk on fire")

	out := buf.String()
	if !strings.Contains(out, "[ERROR]") {
		t.Fatalf("expected output to contain level tag, got %q", out)
	}
	if !strings.Contains(out, "disk on fire") {
		t.Fatalf("expected output to contain message, got %q", out)
	}
}

func TestFileAppenderWritesToTempFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "app.log")

	appender, err := NewFileAppender(path)
	if err != nil {
		t.Fatalf("unexpected err creating file appender: %v", err)
	}

	logger := NewLogger(INFO, appender)
	logger.Info("first line")
	logger.Warn("second line")
	logger.Debug("filtered, should not appear")

	if err := appender.Close(); err != nil {
		t.Fatalf("unexpected err closing file appender: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("unexpected err reading log file: %v", err)
	}
	content := string(data)

	lines := strings.Split(strings.TrimRight(content, "\n"), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines in log file, got %d: %q", len(lines), content)
	}
	if !strings.Contains(lines[0], "first line") {
		t.Fatalf("expected first line to contain %q, got %q", "first line", lines[0])
	}
	if !strings.Contains(lines[1], "second line") {
		t.Fatalf("expected second line to contain %q, got %q", "second line", lines[1])
	}
	if strings.Contains(content, "filtered, should not appear") {
		t.Fatalf("did not expect filtered debug message in file: %q", content)
	}
}

// TestConcurrentLogCalls exercises Log() from many goroutines simultaneously
// to verify the logger and its appenders stay consistent under concurrent
// use (no panics, no lost/corrupted records, no data races under -race).
func TestConcurrentLogCalls(t *testing.T) {
	const goroutines = 100
	const perGoroutine = 20

	appender := &mockAppender{}
	logger := NewLogger(DEBUG, appender)

	var wg sync.WaitGroup
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < perGoroutine; j++ {
				logger.Info("concurrent message")
			}
		}(i)
	}
	wg.Wait()

	want := goroutines * perGoroutine
	if got := appender.count(); got != want {
		t.Fatalf("expected %d records, got %d", want, got)
	}
}
