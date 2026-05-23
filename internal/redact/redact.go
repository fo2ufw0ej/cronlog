// Package redact provides functionality to mask sensitive values
// (e.g. secrets, tokens) found in log output entries.
package redact

import (
	"regexp"
	"strings"

	"github.com/cronlog/internal/logger"
)

const mask = "[REDACTED]"

// Options controls how redaction is applied.
type Options struct {
	// Patterns is a list of literal strings to mask.
	Patterns []string
	// Regexps is a list of regular expression patterns whose matches are masked.
	Regexps []string
}

// Apply returns a new slice of entries with sensitive values replaced by [REDACTED].
// Original entries are not modified.
func Apply(entries []logger.Entry, opts Options) ([]logger.Entry, error) {
	if len(opts.Patterns) == 0 && len(opts.Regexps) == 0 {
		return entries, nil
	}

	compiledRegexps := make([]*regexp.Regexp, 0, len(opts.Regexps))
	for _, pattern := range opts.Regexps {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
		compiledRegexps = append(compiledRegexps, re)
	}

	result := make([]logger.Entry, len(entries))
	for i, e := range entries {
		line := e.Line
		for _, p := range opts.Patterns {
			if p != "" {
				line = strings.ReplaceAll(line, p, mask)
			}
		}
		for _, re := range compiledRegexps {
			line = re.ReplaceAllString(line, mask)
		}
		result[i] = logger.Entry{
			Timestamp: e.Timestamp,
			Stream:    e.Stream,
			Line:      line,
		}
	}
	return result, nil
}
