package logger_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/user/cronlog/internal/logger"
)

func TestWriter_SingleLine(t *testing.T) {
	log := logger.New()
	w := log.Writer("stdout")

	fmt.Fprintln(w, "hello world")

	entries := log.AllEntries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if !strings.Contains(entries[0].Line, "hello world") {
		t.Errorf("unexpected entry line: %q", entries[0].Line)
	}
	if entries[0].Stream != "stdout" {
		t.Errorf("expected stream stdout, got %q", entries[0].Stream)
	}
}

func TestWriter_MultipleLines(t *testing.T) {
	log := logger.New()
	w := log.Writer("stderr")

	fmt.Fprintf(w, "line one\nline two\nline three\n")

	entries := log.AllEntries()
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
}

func TestWriter_MixedStreams(t *testing.T) {
	log := logger.New()
	out := log.Writer("stdout")
	err := log.Writer("stderr")

	fmt.Fprintln(out, "from stdout")
	fmt.Fprintln(err, "from stderr")
	fmt.Fprintln(out, "stdout again")

	entries := log.AllEntries()
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}

	streams := map[string]int{}
	for _, e := range entries {
		streams[e.Stream]++
	}
	if streams["stdout"] != 2 {
		t.Errorf("expected 2 stdout entries, got %d", streams["stdout"])
	}
	if streams["stderr"] != 1 {
		t.Errorf("expected 1 stderr entry, got %d", streams["stderr"])
	}
}
