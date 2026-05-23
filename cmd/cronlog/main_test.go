package main

import (
	"testing"
	"time"

	"github.com/yourorg/cronlog/internal/logger"
	"github.com/yourorg/cronlog/internal/runner"
)

func TestBuildPayload_IncludesFields(t *testing.T) {
	log := logger.New(nil)
	log.Stdout().Write([]byte("hello world\n")) //nolint:errcheck

	result := &runner.Result{
		ExitCode: 1,
		Started:  time.Now().Add(-2 * time.Second),
		Finished: time.Now(),
		Log:      log,
	}

	payload := buildPayload([]string{"echo", "hello"}, result)

	if payload["exit_code"] != 1 {
		t.Errorf("expected exit_code=1, got %v", payload["exit_code"])
	}

	cmd, ok := payload["command"].([]string)
	if !ok || len(cmd) != 2 || cmd[0] != "echo" {
		t.Errorf("unexpected command field: %v", payload["command"])
	}

	output, ok := payload["output"].([]string)
	if !ok {
		t.Fatalf("output field missing or wrong type")
	}
	if len(output) == 0 {
		t.Error("expected at least one output line")
	}

	dur, ok := payload["duration"].(string)
	if !ok || dur == "" {
		t.Errorf("expected non-empty duration string, got %v", payload["duration"])
	}
}

func TestBuildPayload_EmptyOutput(t *testing.T) {
	log := logger.New(nil)

	result := &runner.Result{
		ExitCode: 0,
		Started:  time.Now(),
		Finished: time.Now(),
		Log:      log,
	}

	payload := buildPayload([]string{"true"}, result)

	output, ok := payload["output"].([]string)
	if !ok {
		t.Fatal("output field missing or wrong type")
	}
	if len(output) != 0 {
		t.Errorf("expected empty output, got %v", output)
	}
}
