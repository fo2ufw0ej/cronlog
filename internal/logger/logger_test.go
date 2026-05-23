package logger

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestWrite_StoresEntry(t *testing.T) {
	var buf bytes.Buffer
	l := New(&buf, LevelStdout)

	_, err := l.Write([]byte("hello world\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(l.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(l.Entries))
	}

	entry := l.Entries[0]
	if entry.Level != LevelStdout {
		t.Errorf("expected level %q, got %q", LevelStdout, entry.Level)
	}
	if entry.Message != "hello world\n" {
		t.Errorf("unexpected message: %q", entry.Message)
	}
	if time.Since(entry.Timestamp) > 2*time.Second {
		t.Errorf("timestamp looks stale: %v", entry.Timestamp)
	}
}

func TestWrite_FormatsOutput(t *testing.T) {
	var buf bytes.Buffer
	l := New(&buf, LevelStderr)

	l.Write([]byte("error occurred"))

	output := buf.String()
	if !strings.Contains(output, "[stderr]") {
		t.Errorf("expected '[stderr]' in output, got: %s", output)
	}
	if !strings.Contains(output, "error occurred") {
		t.Errorf("expected message in output, got: %s", output)
	}
}

func TestAllEntries_MergesAndSorts(t *testing.T) {
	var buf bytes.Buffer

	out := New(&buf, LevelStdout)
	errL := New(&buf, LevelStderr)

	// Write to stderr first, then stdout — AllEntries should sort by timestamp.
	errL.Entries = append(errL.Entries, Entry{
		Timestamp: time.Now().Add(-2 * time.Second),
		Level:     LevelStderr,
		Message:   "first",
	})
	out.Entries = append(out.Entries, Entry{
		Timestamp: time.Now(),
		Level:     LevelStdout,
		Message:   "second",
	})

	entries := AllEntries(out, errL)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Message != "first" {
		t.Errorf("expected 'first' to be sorted first, got %q", entries[0].Message)
	}
	if entries[1].Message != "second" {
		t.Errorf("expected 'second' to be sorted second, got %q", entries[1].Message)
	}
}
