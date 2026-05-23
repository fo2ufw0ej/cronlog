package runner_test

import (
	"strings"
	"testing"

	"github.com/user/cronlog/internal/logger"
	"github.com/user/cronlog/internal/runner"
)

func newRunner() (*runner.Runner, *logger.Logger) {
	log := logger.New()
	return runner.New(log), log
}

func TestRun_Success(t *testing.T) {
	r, _ := newRunner()

	result, err := r.Run("echo", "hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success() {
		t.Errorf("expected success, got exit code %d", result.ExitCode)
	}
	if len(result.Entries) == 0 {
		t.Error("expected at least one log entry")
	}
	found := false
	for _, e := range result.Entries {
		if strings.Contains(e.Line, "hello") {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected log entry containing 'hello'")
	}
}

func TestRun_NonZeroExit(t *testing.T) {
	r, _ := newRunner()

	result, err := r.Run("sh", "-c", "exit 2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Success() {
		t.Error("expected failure")
	}
	if result.ExitCode != 2 {
		t.Errorf("expected exit code 2, got %d", result.ExitCode)
	}
}

func TestRun_InvalidCommand(t *testing.T) {
	r, _ := newRunner()

	_, err := r.Run("__nonexistent_command_xyz__")
	if err == nil {
		t.Error("expected error for nonexistent command")
	}
}

func TestResult_Duration(t *testing.T) {
	r, _ := newRunner()

	result, err := r.Run("echo", "timing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Duration() <= 0 {
		t.Error("expected positive duration")
	}
}
