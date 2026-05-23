package summary

import (
	"fmt"
	"time"

	"github.com/example/cronlog/internal/logger"
)

// Stats holds aggregated statistics about a job run.
type Stats struct {
	TotalLines  int
	StdoutLines int
	StderrLines int
	Duration    time.Duration
	ExitCode    int
	StartTime   time.Time
	EndTime     time.Time
}

// Build computes a Stats summary from a slice of log entries and run metadata.
func Build(entries []logger.Entry, exitCode int, start, end time.Time) Stats {
	var stdout, stderr int
	for _, e := range entries {
		switch e.Stream {
		case "stdout":
			stdout++
		case "stderr":
			stderr++
		}
	}
	return Stats{
		TotalLines:  len(entries),
		StdoutLines: stdout,
		StderrLines: stderr,
		Duration:    end.Sub(start),
		ExitCode:    exitCode,
		StartTime:   start,
		EndTime:     end,
	}
}

// Format returns a human-readable one-line summary string.
func (s Stats) Format() string {
	status := "ok"
	if s.ExitCode != 0 {
		status = fmt.Sprintf("exit %d", s.ExitCode)
	}
	return fmt.Sprintf(
		"status=%s duration=%s lines=%d (stdout=%d stderr=%d)",
		status,
		s.Duration.Round(time.Millisecond),
		s.TotalLines,
		s.StdoutLines,
		s.StderrLines,
	)
}
