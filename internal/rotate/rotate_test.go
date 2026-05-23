package rotate_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/example/cronlog/internal/rotate"
)

func TestOpen_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	r := rotate.New(rotate.Options{Dir: dir, Prefix: "job"})
	ts := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

	f, err := r.Open(ts)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer f.Close()

	if _, err := os.Stat(f.Name()); err != nil {
		t.Errorf("expected file to exist: %v", err)
	}
	base := filepath.Base(f.Name())
	expected := "job-20240601T120000Z.log"
	if base != expected {
		t.Errorf("filename = %q, want %q", base, expected)
	}
}

func TestOpen_CreatesDir(t *testing.T) {
	base := t.TempDir()
	dir := filepath.Join(base, "sub", "logs")
	r := rotate.New(rotate.Options{Dir: dir, Prefix: "cron"})

	f, err := r.Open(time.Now())
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	f.Close()

	if _, err := os.Stat(dir); err != nil {
		t.Errorf("expected dir to be created: %v", err)
	}
}

func TestPrune_RemovesOldest(t *testing.T) {
	dir := t.TempDir()
	r := rotate.New(rotate.Options{Dir: dir, Prefix: "job", MaxFiles: 2})

	times := []time.Time{
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
	}
	for _, ts := range times {
		f, err := r.Open(ts)
		if err != nil {
			t.Fatalf("Open: %v", err)
		}
		f.Close()
	}

	if err := r.Prune(); err != nil {
		t.Fatalf("Prune: %v", err)
	}

	matches, _ := filepath.Glob(filepath.Join(dir, "*.log"))
	if len(matches) != 2 {
		t.Errorf("after prune: %d files, want 2", len(matches))
	}
}

func TestPrune_NoLimit_KeepsAll(t *testing.T) {
	dir := t.TempDir()
	r := rotate.New(rotate.Options{Dir: dir, Prefix: "job", MaxFiles: 0})

	for i := 0; i < 5; i++ {
		ts := time.Date(2024, 1, i+1, 0, 0, 0, 0, time.UTC)
		f, _ := r.Open(ts)
		f.Close()
	}

	if err := r.Prune(); err != nil {
		t.Fatalf("Prune: %v", err)
	}

	matches, _ := filepath.Glob(filepath.Join(dir, "*.log"))
	if len(matches) != 5 {
		t.Errorf("after prune: %d files, want 5", len(matches))
	}
}
