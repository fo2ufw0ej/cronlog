package webhook_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/cronlog/internal/webhook"
)

func TestSend_Success(t *testing.T) {
	var received webhook.Payload

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected application/json, got %s", ct)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := webhook.New(server.URL)
	payload := webhook.Payload{
		Job:       "backup",
		StartedAt: time.Now(),
		Duration:  3.5,
		ExitCode:  0,
		Output:    "done",
		Success:   true,
	}

	if err := client.Send(payload); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received.Job != "backup" {
		t.Errorf("expected job 'backup', got %q", received.Job)
	}
	if !received.Success {
		t.Error("expected success to be true")
	}
}

func TestSend_NonSuccessStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := webhook.New(server.URL)
	err := client.Send(webhook.Payload{Job: "test"})
	if err == nil {
		t.Fatal("expected error for non-2xx status, got nil")
	}
}

func TestSend_InvalidURL(t *testing.T) {
	client := webhook.New("http://127.0.0.1:0/invalid")
	err := client.Send(webhook.Payload{Job: "test"})
	if err == nil {
		t.Fatal("expected error for unreachable URL, got nil")
	}
}
