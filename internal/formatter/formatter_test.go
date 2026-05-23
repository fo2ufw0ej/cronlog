package formatter_test

import (
	"strings"
	"testing"
	"time"

	"github.com/example/cronlog/internal/formatter"
	"github.com/example/cronlog/internal/logger"
)

func makeEntries() []logger.Entry {
	t1 := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, 1, 15, 10, 0, 1, 0, time.UTC)
	return []logger.Entry{
		{Timestamp: t1, Stream: "stdout", Message: "hello world"},
		{Timestamp: t2, Stream: "stderr", Message: "an error occurred"},
	}
}

func TestRender_TextFormat_WithStream(t *testing.T) {
	entries := makeEntries()
	opts := formatter.DefaultOptions()
	out, err := formatter.Render(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "[stdout]") {
		t.Error("expected stdout stream label in output")
	}
	if !strings.Contains(out, "hello world") {
		t.Error("expected message in output")
	}
	if !strings.Contains(out, "2024-01-15T10:00:00Z") {
		t.Error("expected timestamp in output")
	}
}

func TestRender_TextFormat_HideStream(t *testing.T) {
	entries := makeEntries()
	opts := formatter.DefaultOptions()
	opts.ShowStream = false
	out, err := formatter.Render(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "[stdout]") {
		t.Error("expected stream label to be hidden")
	}
}

func TestRender_JSONFormat(t *testing.T) {
	entries := makeEntries()
	opts := formatter.DefaultOptions()
	opts.Format = formatter.FormatJSON
	out, err := formatter.Render(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"message"`) {
		t.Error("expected JSON message field")
	}
	if !strings.Contains(out, `"timestamp"`) {
		t.Error("expected JSON timestamp field")
	}
}

func TestRender_InvalidTimezone(t *testing.T) {
	entries := makeEntries()
	opts := formatter.DefaultOptions()
	opts.TimeZone = "Not/AReal_Zone"
	_, err := formatter.Render(entries, opts)
	if err == nil {
		t.Error("expected error for invalid timezone")
	}
}

func TestRender_EmptyEntries(t *testing.T) {
	out, err := formatter.Render([]logger.Entry{}, formatter.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "" {
		t.Errorf("expected empty output, got %q", out)
	}
}
