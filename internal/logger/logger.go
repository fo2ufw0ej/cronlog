package logger

import (
	"fmt"
	"io"
	"time"
)

// Level represents the log level for a message.
type Level string

const (
	LevelStdout Level = "stdout"
	LevelStderr Level = "stderr"
	LevelInfo   Level = "info"
)

// Entry represents a single timestamped log entry.
type Entry struct {
	Timestamp time.Time
	Level     Level
	Message   string
}

// Logger wraps an io.Writer and prefixes each line with a timestamp and level.
type Logger struct {
	w       io.Writer
	level   Level
	Entries []Entry
}

// New creates a new Logger that writes to w with the given level.
func New(w io.Writer, level Level) *Logger {
	return &Logger{
		w:     w,
		level: level,
	}
}

// Write implements io.Writer. Each call is treated as a log entry.
func (l *Logger) Write(p []byte) (n int, err error) {
	now := time.Now().UTC()
	msg := string(p)

	entry := Entry{
		Timestamp: now,
		Level:     l.level,
		Message:   msg,
	}
	l.Entries = append(l.Entries, entry)

	line := fmt.Sprintf("[%s] [%s] %s", now.Format(time.RFC3339), l.level, msg)
	_, err = fmt.Fprint(l.w, line)
	return len(p), err
}

// AllEntries merges and returns entries from multiple loggers sorted by timestamp.
func AllEntries(loggers ...*Logger) []Entry {
	var all []Entry
	for _, l := range loggers {
		all = append(all, l.Entries...)
	}
	// Simple insertion sort — log volumes are typically small.
	for i := 1; i < len(all); i++ {
		for j := i; j > 0 && all[j].Timestamp.Before(all[j-1].Timestamp); j-- {
			all[j], all[j-1] = all[j-1], all[j]
		}
	}
	return all
}
