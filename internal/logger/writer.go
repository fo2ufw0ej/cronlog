package logger

import (
	"bufio"
	"io"
	"strings"
)

// streamWriter is an io.Writer that feeds lines into the Logger
// tagged with a given stream name (e.g. "stdout" or "stderr").
type streamWriter struct {
	logger *Logger
	stream string
	buf    strings.Builder
}

// Writer returns an io.Writer that logs each line written to it
// under the given stream label.
func (l *Logger) Writer(stream string) io.Writer {
	return &streamWriter{logger: l, stream: stream}
}

// Write implements io.Writer. It buffers data and flushes complete lines
// to the underlying logger.
func (sw *streamWriter) Write(p []byte) (int, error) {
	sw.buf.Write(p)
	scanner := bufio.NewScanner(strings.NewReader(sw.buf.String()))
	var remaining strings.Builder
	for scanner.Scan() {
		line := scanner.Text()
		sw.logger.Write(sw.stream, line) //nolint:errcheck
	}
	// If the last byte is not a newline, keep the partial line buffered.
	raw := sw.buf.String()
	if len(raw) > 0 && raw[len(raw)-1] != '\n' {
		parts := strings.Split(raw, "\n")
		remaining.WriteString(parts[len(parts)-1])
	}
	sw.buf.Reset()
	sw.buf.WriteString(remaining.String())
	return len(p), nil
}
