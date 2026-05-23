package summary_test

import (
	"testing"
	"time"

	"github.com/example/cronlog/internal/logger"
	"github.com/example/cronlog/internal/summary"
)

func makeEntries(stdout, stderr int) []logger.Entry {
	var entries []logger.Entry
	now := time.Now()
	for i := 0; i < stdout; i++ {
		entries = append(entries, logger.Entry{Stream: "stdout", Line: "out", Timestamp: now})
	}
	for i := 0; i < stderr; i++ {
		entries = append(entries, logger.Entry{Stream: "stderr", Line: "err", Timestamp: now})
	}
	return entries
}

func TestBuild_CountsLines(t *testing.T) {
	start := time.Now()
	end := start.Add(2 * time.Second)
	entries := makeEntries(3, 2)

	s := summary.Build(entries, 0, start, end)

	if s.TotalLines != 5 {
		t.Errorf("TotalLines: got %d, want 5", s.TotalLines)
	}
	if s.StdoutLines != 3 {
		t.Errorf("StdoutLines: got %d, want 3", s.StdoutLines)
	}
	if s.StderrLines != 2 {
		t.Errorf("StderrLines: got %d, want 2", s.StderrLines)
	}
}

func TestBuild_Duration(t *testing.T) {
	start := time.Now()
	end := start.Add(500 * time.Millisecond)

	s := summary.Build(nil, 0, start, end)

	if s.Duration != 500*time.Millisecond {
		t.Errorf("Duration: got %v, want 500ms", s.Duration)
	}
}

func TestFormat_SuccessStatus(t *testing.T) {
	start := time.Now()
	end := start.Add(100 * time.Millisecond)
	s := summary.Build(makeEntries(1, 0), 0, start, end)

	out := s.Format()
	if out == "" {
		t.Fatal("Format returned empty string")
	}
	if !contains(out, "status=ok") {
		t.Errorf("expected status=ok in %q", out)
	}
}

func TestFormat_NonZeroExit(t *testing.T) {
	start := time.Now()
	end := start.Add(50 * time.Millisecond)
	s := summary.Build(nil, 2, start, end)

	out := s.Format()
	if !contains(out, "exit 2") {
		t.Errorf("expected 'exit 2' in %q", out)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
