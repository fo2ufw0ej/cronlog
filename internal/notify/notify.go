package notify

import (
	"fmt"

	"github.com/user/cronlog/internal/formatter"
	"github.com/user/cronlog/internal/logger"
	"github.com/user/cronlog/internal/webhook"
)

// Condition controls when a notification is sent.
type Condition string

const (
	ConditionAlways  Condition = "always"
	ConditionFailure Condition = "failure"
	ConditionNever   Condition = "never"
)

// Options holds notification configuration.
type Options struct {
	WebhookURL string
	Condition  Condition
	FmtOptions formatter.Options
}

// Notifier sends webhook notifications based on job results.
type Notifier struct {
	opts Options
	hook *webhook.Webhook
}

// New creates a Notifier with the given options.
func New(opts Options) *Notifier {
	var hook *webhook.Webhook
	if opts.WebhookURL != "" {
		hook = webhook.New(opts.WebhookURL)
	}
	return &Notifier{opts: opts, hook: hook}
}

// Notify sends a notification if the condition is met.
// failed indicates whether the job exited with a non-zero status.
func (n *Notifier) Notify(jobName string, failed bool, entries []logger.Entry) error {
	if n.hook == nil {
		return nil
	}
	switch n.opts.Condition {
	case ConditionNever:
		return nil
	case ConditionFailure:
		if !failed {
			return nil
		}
	}

	body := formatter.Render(entries, n.opts.FmtOptions)
	payload := map[string]string{
		"job":    jobName,
		"status": statusLabel(failed),
		"output": body,
	}
	return n.hook.Send(payload)
}

func statusLabel(failed bool) string {
	if failed {
		return "failure"
	}
	return fmt.Sprintf("success")
}
