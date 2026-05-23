package filter_test

import (
	"testing"
	"time"

	"github.com/user/cronlog/internal/filter"
	"github.com/user/cronlog/internal/logger"
)

func makeEntries() []logger.Entry {
	now := time.Now()
	return []logger.Entry{
		{Stream: logger.Stdout, Text: "hello world", At: now},
		{Stream: logger.Stderr, Text: "error occurred", At: now.Add(time.Second)},
		{Stream: logger.Stdout, Text: "goodbye world", At: now.Add(2 * time.Second)},
		{Stream: logger.Stderr, Text: "another error", At: now.Add(3 * time.Second)},
	}
}

func TestApply_NoFilter_ReturnsAll(t *testing.T) {
	entries := makeEntries()
	got, err := filter.Apply(entries, filter.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != len(entries) {
		t.Errorf("expected %d entries, got %d", len(entries), len(got))
	}
}

func TestApply_StreamFilter_Stdout(t *testing.T) {
	got, err := filter.Apply(makeEntries(), filter.Options{Stream: "stdout"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 stdout entries, got %d", len(got))
	}
	for _, e := range got {
		if e.Stream != logger.Stdout {
			t.Errorf("expected stdout, got %s", e.Stream)
		}
	}
}

func TestApply_PatternFilter_Matches(t *testing.T) {
	got, err := filter.Apply(makeEntries(), filter.Options{Pattern: "error"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 matching entries, got %d", len(got))
	}
}

func TestApply_CombinedFilter(t *testing.T) {
	got, err := filter.Apply(makeEntries(), filter.Options{Stream: "stderr", Pattern: "another"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("expected 1 entry, got %d", len(got))
	}
}

func TestApply_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := filter.Apply(makeEntries(), filter.Options{Pattern: "[invalid"})
	if err == nil {
		t.Error("expected error for invalid regex, got nil")
	}
}
