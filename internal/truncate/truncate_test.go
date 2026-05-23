package truncate_test

import (
	"testing"
	"time"

	"github.com/user/cronlog/internal/logger"
	"github.com/user/cronlog/internal/truncate"
)

func makeEntries(messages ...string) []logger.Entry {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	entries := make([]logger.Entry, len(messages))
	for i, m := range messages {
		entries[i] = logger.Entry{
			Timestamp: base.Add(time.Duration(i) * time.Second),
			Stream:    "stdout",
			Message:   m,
		}
	}
	return entries
}

func TestApply_NoLimits_ReturnsAll(t *testing.T) {
	entries := makeEntries("a", "b", "c")
	result := truncate.Apply(entries, truncate.Options{})
	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
}

func TestApply_MaxLines_TruncatesAndAppendsSummary(t *testing.T) {
	entries := makeEntries("line1", "line2", "line3", "line4", "line5")
	result := truncate.Apply(entries, truncate.Options{MaxLines: 3})
	// 3 kept + 1 summary
	if len(result) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(result))
	}
	last := result[len(result)-1]
	if last.Stream != "stderr" {
		t.Errorf("summary stream: got %q, want %q", last.Stream, "stderr")
	}
	if last.Message == "" {
		t.Error("summary message should not be empty")
	}
}

func TestApply_MaxBytes_TruncatesAndAppendsSummary(t *testing.T) {
	// Each message is 5 bytes; limit to 12 bytes → 2 entries kept
	entries := makeEntries("hello", "world", "extra")
	result := truncate.Apply(entries, truncate.Options{MaxBytes: 12})
	if len(result) != 3 { // 2 kept + 1 summary
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
}

func TestApply_Empty_ReturnsEmpty(t *testing.T) {
	result := truncate.Apply([]logger.Entry{}, truncate.Options{MaxLines: 5})
	if len(result) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(result))
	}
}

func TestApply_WithinLimits_NoSummaryAdded(t *testing.T) {
	entries := makeEntries("a", "b")
	result := truncate.Apply(entries, truncate.Options{MaxLines: 10, MaxBytes: 1000})
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	for _, e := range result {
		if e.Stream == "stderr" {
			t.Error("unexpected summary entry added when within limits")
		}
	}
}
