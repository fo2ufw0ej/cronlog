package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/cronlog/internal/config"
	"github.com/yourorg/cronlog/internal/logger"
	"github.com/yourorg/cronlog/internal/runner"
	"github.com/yourorg/cronlog/internal/webhook"
)

func main() {
	configPath := flag.String("config", "", "path to config file (optional)")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: cronlog [--config FILE] COMMAND [ARGS...]")
		os.Exit(1)
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cronlog: failed to load config: %v\n", err)
		os.Exit(1)
	}

	log := logger.New(os.Stdout)

	r := runner.New(args[0], args[1:]...)
	result, err := r.Run(log)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cronlog: failed to run command: %v\n", err)
		os.Exit(1)
	}

	if cfg.WebhookURL != "" && (result.ExitCode != 0 || cfg.WebhookAlways) {
		wh := webhook.New(cfg.WebhookURL)
		payload := buildPayload(args, result)
		if whErr := wh.Send(payload); whErr != nil {
			fmt.Fprintf(os.Stderr, "cronlog: webhook error: %v\n", whErr)
		}
	}

	os.Exit(result.ExitCode)
}

func buildPayload(args []string, result *runner.Result) map[string]interface{} {
	entries := result.Log.AllEntries()
	lines := make([]string, 0, len(entries))
	for _, e := range entries {
		lines = append(lines, fmt.Sprintf("[%s] %s: %s", e.Time.Format("15:04:05"), e.Stream, e.Line))
	}
	return map[string]interface{}{
		"command":   args,
		"exit_code": result.ExitCode,
		"duration":  result.Duration().String(),
		"output":    lines,
	}
}
