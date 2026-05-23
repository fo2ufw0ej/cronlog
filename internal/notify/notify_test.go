package notify_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/cronlog/internal/formatter"
	"github.com/user/cronlog/internal/logger"
	"github.com/user/cronlog/internal/notify"
)

func makeEntries() []logger.Entry {
	return []logger.Entry{
		{Timestamp: time.Now(), Stream: "stdout", Message: "hello"},
		{Timestamp: time.Now(), Stream: "stderr", Message: "warn"},
	}
}

func TestNotify_AlwaysCondition_SendsOnSuccess(t *testing.T) {
	var received map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notify.New(notify.Options{
		WebhookURL: ts.URL,
		Condition:  notify.ConditionAlways,
		FmtOptions: formatter.DefaultOptions(),
	})
	if err := n.Notify("myjob", false, makeEntries()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["job"] != "myjob" {
		t.Errorf("expected job=myjob, got %q", received["job"])
	}
	if received["status"] != "success" {
		t.Errorf("expected status=success, got %q", received["status"])
	}
}

func TestNotify_FailureCondition_SkipsOnSuccess(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notify.New(notify.Options{
		WebhookURL: ts.URL,
		Condition:  notify.ConditionFailure,
		FmtOptions: formatter.DefaultOptions(),
	})
	if err := n.Notify("myjob", false, makeEntries()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected webhook NOT to be called on success with failure condition")
	}
}

func TestNotify_FailureCondition_SendsOnFailure(t *testing.T) {
	var received map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notify.New(notify.Options{
		WebhookURL: ts.URL,
		Condition:  notify.ConditionFailure,
		FmtOptions: formatter.DefaultOptions(),
	})
	if err := n.Notify("myjob", true, makeEntries()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["status"] != "failure" {
		t.Errorf("expected status=failure, got %q", received["status"])
	}
}

func TestNotify_NeverCondition_NeverSends(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notify.New(notify.Options{
		WebhookURL: ts.URL,
		Condition:  notify.ConditionNever,
		FmtOptions: formatter.DefaultOptions(),
	})
	_ = n.Notify("myjob", true, makeEntries())
	if called {
		t.Error("expected webhook NOT to be called with never condition")
	}
}

func TestNotify_NoWebhookURL_IsNoop(t *testing.T) {
	n := notify.New(notify.Options{
		Condition:  notify.ConditionAlways,
		FmtOptions: formatter.DefaultOptions(),
	})
	if err := n.Notify("myjob", true, makeEntries()); err != nil {
		t.Fatalf("expected no error with empty webhook URL, got: %v", err)
	}
}
