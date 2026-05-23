package runner

import (
	"os/exec"
	"time"

	"github.com/user/cronlog/internal/logger"
)

// Result holds the outcome of a command execution.
type Result struct {
	Command  string
	Args     []string
	ExitCode int
	Started  time.Time
	Finished time.Time
	Entries  []logger.Entry
}

// Runner executes commands and captures their output via a logger.
type Runner struct {
	log *logger.Logger
}

// New creates a new Runner backed by the provided logger.
func New(log *logger.Logger) *Runner {
	return &Runner{log: log}
}

// Run executes the given command with args, streaming stdout and stderr
// through the logger. It returns a Result regardless of exit status.
// If the command fails to start (e.g. binary not found), an error is returned
// along with a partial Result containing timing and any collected entries.
func (r *Runner) Run(command string, args ...string) (*Result, error) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = r.log.Writer("stdout")
	cmd.Stderr = r.log.Writer("stderr")

	result := &Result{
		Command: command,
		Args:    args,
		Started: time.Now(),
	}

	err := cmd.Run()
	result.Finished = time.Now()
	result.Entries = r.log.AllEntries()

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
			return result, nil
		}
		return result, err
	}

	return result, nil
}

// Duration returns how long the command took to run.
func (r *Result) Duration() time.Duration {
	return r.Finished.Sub(r.Started)
}

// Success returns true if the command exited with code 0.
func (r *Result) Success() bool {
	return r.ExitCode == 0
}

// StderrEntries returns only the log entries captured from stderr.
func (r *Result) StderrEntries() []logger.Entry {
	var entries []logger.Entry
	for _, e := range r.Entries {
		if e.Stream == "stderr" {
			entries = append(entries, e)
		}
	}
	return entries
}
