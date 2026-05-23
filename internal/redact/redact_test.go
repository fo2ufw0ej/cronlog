package redact_test

import (
	"testing"
	"time"

	"github.com/cronlog/internal/logger"
	"github.com/cronlog/internal/redact"
)

func makeEntries(lines ...string) []logger.Entry {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	entries := make([]logger.Entry, len(lines))
	for i, l := range lines {
		entries[i] = logger.Entry{
			Timestamp: base.Add(time.Duration(i) * time.Second),
			Stream:    "stdout",
			Line:      l,
		}
	}
	return entries
}

func TestApply_NoOptions_ReturnsOriginal(t *testing.T) {
	entries := makeEntries("hello world", "foo bar")
	out, err := redact.Apply(entries, redact.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != len(entries) {
		t.Fatalf("expected %d entries, got %d", len(entries), len(out))
	}
	for i, e := range out {
		if e.Line != entries[i].Line {
			t.Errorf("entry %d: expected %q, got %q", i, entries[i].Line, e.Line)
		}
	}
}

func TestApply_LiteralPattern_Masked(t *testing.T) {
	entries := makeEntries("token=supersecret value", "nothing here")
	out, err := redact.Apply(entries, redact.Options{Patterns: []string{"supersecret"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Line != "token=[REDACTED] value" {
		t.Errorf("expected redacted line, got %q", out[0].Line)
	}
	if out[1].Line != "nothing here" {
		t.Errorf("expected unchanged line, got %q", out[1].Line)
	}
}

func TestApply_RegexpPattern_Masked(t *testing.T) {
	entries := makeEntries("password=abc123", "user=admin")
	out, err := redact.Apply(entries, redact.Options{Regexps: []string{`password=\S+`}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Line != "[REDACTED]" {
		t.Errorf("expected redacted line, got %q", out[0].Line)
	}
	if out[1].Line != "user=admin" {
		t.Errorf("expected unchanged line, got %q", out[1].Line)
	}
}

func TestApply_InvalidRegexp_ReturnsError(t *testing.T) {
	entries := makeEntries("some line")
	_, err := redact.Apply(entries, redact.Options{Regexps: []string{`[invalid(`}})
	if err == nil {
		t.Fatal("expected error for invalid regexp, got nil")
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	entries := makeEntries("secret=mysecret")
	original := entries[0].Line
	_, err := redact.Apply(entries, redact.Options{Patterns: []string{"mysecret"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries[0].Line != original {
		t.Errorf("original entry was mutated: got %q", entries[0].Line)
	}
}
