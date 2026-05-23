# cronlog

Wraps cron jobs to capture stdout/stderr with timestamps and optional webhook notifications.

## Installation

```bash
go install github.com/yourusername/cronlog@latest
```

## Usage

Replace your cron command with `cronlog` as a wrapper:

```
# crontab -e
0 2 * * * cronlog --tag "nightly-backup" /usr/local/bin/backup.sh
```

cronlog captures all output, prepends timestamps to each line, and writes to a log file.

### Options

| Flag | Description | Default |
|------|-------------|---------|
| `--tag` | Label for this job | command name |
| `--log` | Log file path | `/var/log/cronlog/<tag>.log` |
| `--webhook` | URL to POST results on failure | _(none)_ |
| `--on-success` | Also notify webhook on success | `false` |

### Example Output

```
2024-03-15T02:00:01Z [nightly-backup] Starting job: /usr/local/bin/backup.sh
2024-03-15T02:00:03Z [nightly-backup] Backup complete: 3 files written
2024-03-15T02:00:03Z [nightly-backup] Job exited with code 0 in 2.1s
```

### Webhook Payload

When a webhook URL is provided, cronlog sends a JSON POST on job completion:

```json
{
  "tag": "nightly-backup",
  "exit_code": 1,
  "duration": "2.1s",
  "output": "..."
}
```

## License

MIT