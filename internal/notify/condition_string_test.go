package notify_test

import (
	"testing"

	"github.com/user/cronlog/internal/notify"
)

func TestParseCondition_ValidValues(t *testing.T) {
	cases := []struct {
		input string
		want  notify.Condition
	}{
		{"always", notify.ConditionAlways},
		{"failure", notify.ConditionFailure},
		{"never", notify.ConditionNever},
		{"", notify.ConditionAlways},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, err := notify.ParseCondition(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("ParseCondition(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestParseCondition_InvalidValue_ReturnsError(t *testing.T) {
	_, err := notify.ParseCondition("weekly")
	if err == nil {
		t.Fatal("expected error for unknown condition, got nil")
	}
}

func TestCondition_String(t *testing.T) {
	if notify.ConditionFailure.String() != "failure" {
		t.Errorf("String() = %q, want %q", notify.ConditionFailure.String(), "failure")
	}
}
