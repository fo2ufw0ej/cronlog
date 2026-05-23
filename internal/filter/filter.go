package filter

import (
	"regexp"
	"strings"

	"github.com/user/cronlog/internal/logger"
)

// Options holds filtering criteria for log entries.
type Options struct {
	// Stream filters by stream type: "stdout", "stderr", or "" for all.
	Stream string
	// Pattern is an optional regex to match against entry text.
	Pattern string
}

// Apply returns a subset of entries matching the given options.
func Apply(entries []logger.Entry, opts Options) ([]logger.Entry, error) {
	var re *regexp.Regexp
	if opts.Pattern != "" {
		var err error
		re, err = regexp.Compile(opts.Pattern)
		if err != nil {
			return nil, err
		}
	}

	result := make([]logger.Entry, 0, len(entries))
	for _, e := range entries {
		if opts.Stream != "" && !strings.EqualFold(string(e.Stream), opts.Stream) {
			continue
		}
		if re != nil && !re.MatchString(e.Text) {
			continue
		}
		result = append(result, e)
	}
	return result, nil
}
