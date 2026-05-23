package rotate

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Options configures log rotation behaviour.
type Options struct {
	// Dir is the directory where log files are stored.
	Dir string
	// Prefix is prepended to every log file name.
	Prefix string
	// MaxFiles is the maximum number of log files to retain (0 = unlimited).
	MaxFiles int
}

// Rotator manages creation and pruning of log files.
type Rotator struct {
	opts Options
}

// New returns a Rotator configured with opts.
func New(opts Options) *Rotator {
	return &Rotator{opts: opts}
}

// Open creates a new timestamped log file in opts.Dir and returns it.
// The caller is responsible for closing the file.
func (r *Rotator) Open(t time.Time) (*os.File, error) {
	if err := os.MkdirAll(r.opts.Dir, 0o755); err != nil {
		return nil, fmt.Errorf("rotate: mkdir %s: %w", r.opts.Dir, err)
	}
	name := filepath.Join(r.opts.Dir, r.filename(t))
	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, fmt.Errorf("rotate: open %s: %w", name, err)
	}
	return f, nil
}

// Prune removes the oldest log files so that at most opts.MaxFiles remain.
// If MaxFiles is 0 no files are removed.
func (r *Rotator) Prune() error {
	if r.opts.MaxFiles <= 0 {
		return nil
	}
	entries, err := r.list()
	if err != nil {
		return err
	}
	for len(entries) > r.opts.MaxFiles {
		if err := os.Remove(entries[0]); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("rotate: prune %s: %w", entries[0], err)
		}
		entries = entries[1:]
	}
	return nil
}

// list returns log files in opts.Dir matching the prefix, sorted oldest first.
func (r *Rotator) list() ([]string, error) {
	glob := filepath.Join(r.opts.Dir, r.opts.Prefix+"*.log")
	matches, err := filepath.Glob(glob)
	if err != nil {
		return nil, fmt.Errorf("rotate: glob: %w", err)
	}
	sort.Strings(matches)
	return matches, nil
}

func (r *Rotator) filename(t time.Time) string {
	ts := t.UTC().Format("20060102T150405Z")
	prefix := r.opts.Prefix
	if prefix != "" && !strings.HasSuffix(prefix, "-") {
		prefix += "-"
	}
	return fmt.Sprintf("%s%s.log", prefix, ts)
}
