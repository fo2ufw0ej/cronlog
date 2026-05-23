package truncate

import (
	"fmt"

	"github.com/user/cronlog/internal/logger"
)

// Options controls how output is truncated.
type Options struct {
	// MaxLines is the maximum number of log entries to keep.
	// If zero, no line-based truncation is applied.
	MaxLines int

	// MaxBytes is the maximum total byte length of all entry messages combined.
	// If zero, no byte-based truncation is applied.
	MaxBytes int
}

// Apply truncates entries according to opts, returning the (possibly shortened)
// slice and a summary line appended when truncation occurred.
func Apply(entries []logger.Entry, opts Options) []logger.Entry {
	if len(entries) == 0 {
		return entries
	}

	result := entries
	truncated := false

	// Byte-based truncation first so line limit is applied to the reduced set.
	if opts.MaxBytes > 0 {
		total := 0
		for i, e := range result {
			total += len(e.Message)
			if total > opts.MaxBytes {
				result = result[:i]
				truncated = true
				break
			}
		}
	}

	// Line-based truncation.
	if opts.MaxLines > 0 && len(result) > opts.MaxLines {
		result = result[:opts.MaxLines]
		truncated = true
	}

	if truncated {
		var ts interface{}
		if len(result) > 0 {
			ts = result[len(result)-1].Timestamp
		}
		_ = ts
		summary := logger.Entry{
			Stream:  "stderr",
			Message: fmt.Sprintf("[cronlog: output truncated, %d lines shown]", len(result)),
		}
		if len(result) > 0 {
			summary.Timestamp = result[len(result)-1].Timestamp
		}
		result = append(result, summary)
	}

	return result
}
