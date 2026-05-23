package formatter

import (
	"fmt"
	"strings"
	"time"

	"github.com/example/cronlog/internal/logger"
)

// Format controls the output format of log entries.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Options holds formatting configuration.
type Options struct {
	Format    Format
	TimeZone  string
	ShowStream bool
}

// DefaultOptions returns sensible formatting defaults.
func DefaultOptions() Options {
	return Options{
		Format:    FormatText,
		TimeZone:  "UTC",
		ShowStream: true,
	}
}

// Render formats a slice of log entries according to the given options.
func Render(entries []logger.Entry, opts Options) (string, error) {
	loc, err := time.LoadLocation(opts.TimeZone)
	if err != nil {
		return "", fmt.Errorf("invalid timezone %q: %w", opts.TimeZone, err)
	}

	var sb strings.Builder
	for _, e := range entries {
		ts := e.Timestamp.In(loc).Format(time.RFC3339)
		switch opts.Format {
		case FormatJSON:
			stream := ""
			if opts.ShowStream {
				stream = fmt.Sprintf(`"stream":%q,`, e.Stream)
			}
			sb.WriteString(fmt.Sprintf(`{%s"timestamp":%q,"message":%q}`+"\n",
				stream, ts, e.Message))
		default:
			if opts.ShowStream {
				sb.WriteString(fmt.Sprintf("[%s] [%s] %s\n", ts, e.Stream, e.Message))
			} else {
				sb.WriteString(fmt.Sprintf("[%s] %s\n", ts, e.Message))
			}
		}
	}
	return sb.String(), nil
}
