<img src="images/green-orb.png" align="right" width="150" height="150" />

# Green Orb
> An 'Observe and Report Buddy' for your SRE toolbox

## Introduction

Green Orb is a lightweight monitoring tool that enhances your
application's reliability by observing its console output for specific
patterns and executing predefined actions in response. Designed to
integrate seamlessly, it's deployed as a single executable binary that
runs your application as a subprocess, where it can monitor all
console output, making it particularly useful in containerized
environments. Green Orb acts as a proactive assistant, handling
essential monitoring tasks and enabling SREs to automate responses to
critical system events effectively.

## Features

With Green Orb, you can:

- **Send Notifications**: Utilize popular messaging services like Slack, Email, Discord, and more to keep your team updated.
- **Trigger Webhooks**: Integrate with services like Ansible Automation Platform through webhooks for seamless automation.
- **Publish to Kafka**: Send important alerts or logs directly to a Kafka topic for real-time processing.
- **Execute Commands**: Run shell commands automatically, allowing actions like capturing thread dumps of the observed process.
- **Manage Processes**: Restart or kill the observed process to maintain desired state or recover from issues.
- **Export Metrics**: Expose Prometheus metrics for observability and alerting.
- **Environment Management**: Automatically load `.env` files to configure your application environment securely.

## CLI Flags

- `-c, --config string` path to the configuration file (default `green-orb.yaml`).
- `-w, --workers int` number of reporting workers (default `5`).
- `--metrics-enable` enable Prometheus metrics endpoint (default `false`).
- `--metrics-addr string` metrics listen address (default `127.0.0.1:9090`).
- `--env string` load environment variables from specified file.
- `--skip-dotenv` do not automatically load .env file (default: loads .env if present).

## Quick Start

Green Orb is easy to configure and use. It's distributed as a single binary, `orb`, and requires just one YAML configuration file.

Simply preface your application with `orb` and customize your
`green-orb.yaml` file to enable its special powers.  For example,
instead of:
```
$ java -jar mywebapp.jar
```
...use...
```
$ orb java -jar mywebapp.jar
```

Or, if you are using containers, change...
```
ENTRYPOINT [ "java", "-jar", "jar-file-name.jar" ]
```
...in your Dockerfile, to...
```
ENTRYPOINT [ "orb", "java", "-jar", "jar-file-name.jar" ]
```

This assumes `green-orb.yaml` is in the current directory. Use the
`-c` flag to point it elsewhere. You can pass flags to your command
without special separators; orb stops parsing at the first non-flag.

### Environment Variables

Green Orb automatically loads environment variables from a `.env` file in the current directory if present. This makes it easy to configure your application environment without exposing secrets in command lines or configuration files.

- `.env` files are loaded automatically (no flag required)
- Use `--skip-dotenv` to prevent loading `.env`
- Use `--env myfile.env` to load additional environment files
- Environment variables are available to both your observed process and orb's templates

The config file tells `orb` what to watch for, and what to do.  For
instance, if it contains the following, you'll get an email every time
your application starts up, and a thread dump every time you get a
thread pool exhausted warning.

```
channels:
  - name: startup-email
    type: notify
    url:  "smtp://MYEMAIL@gmail.com:MYPASSWORD@smtp.gmail.com:587/?from=MYEMAIL@gmail.com&to=MYEMAIL@gmail.com&subject=Application%20Starting!"

  - name: thread-dump
    type: exec
    shell: |
      FILENAME=thread-dump-$(date).txt
      jstack $ORB_PID > /tmp/${FILENAME}
      aws s3 mv /tmp/${FILENAME} s3:/my-bucket/${FILENAME}

signals:
  - regex: Starting Application
    channel: startup-email

  - regex: "Warning: thread pool exhausted"
    channel: thread-dump

# Run a periodic action every 5 minutes as a schedule-type signal
signals:
  - name: periodic-stacktrace
    channel: thread-dump
    schedule:
      every: 5m

# Or, run at the top of every hour using cron
signals:
  - name: hourly-stacktrace
    channel: thread-dump
    schedule:
      cron: "0 * * * *"

## Enabling/Disabling Signals at Runtime

Use `enable_signal` and `disable_signal` channels to control a named signal (regex or schedule) dynamically. Signals are enabled by default; you can set `enabled: false` on a signal to start disabled.

Example: enable a schedule signal for 1 hour when a pattern is observed, and provide a separate disable:

```
channels:
  # Turns on the schedule signal for 1 hour
  - name: enable-stacktrace-1h
    type: enable_signal
    target_signal: periodic-stacktrace
    duration: 1h

  # Turns off the schedule signal immediately
  - name: disable-stacktrace
    type: disable_signal
    target_signal: periodic-stacktrace

signals:
  # The schedule-style signal we want to control
  - name: periodic-stacktrace
    channel: thread-dump
    schedule:
      every: 5m

  # When this regex matches, enable periodic stack traces for 1h
  - regex: Enable stack traces
    channel: enable-stacktrace-1h

  # Optional: log pattern to stop early
  - regex: Disable stack traces
    channel: disable-stacktrace
```

Notes:
- `target_signal` must match the `name` of the signal to control.
- `duration` is optional and only applies to enable; after it elapses, the signal auto-disables.
- These fields support templates, so you can compute names or durations from matches or env if needed.
- To start a signal disabled at boot, add `enabled: false` to that signal definition.

Notes:
- Use `schedule.every` for fixed intervals (Go duration strings like `30s`, `5m`, `1h`).
- Use `schedule.cron` for cron-like schedules (5 fields; seconds optional). Cron runs in the system timezone.
- Provide exactly one of `every` or `cron` per signal.

`orb` does not interfere with the execution of your application.  All
console logs still go to the console, Linux and macOS signals
are passed through to the observed process, and the exit code for your
application is returned through `orb`.

## Channels and Signals

In Green Orb, "channels", "signals", and "schedules" are foundational concepts:

- Signals: Mappings of regular expressions to channels. When a
  signal's regex matches a log entry, the corresponding channel is
  invoked.

- Schedules: Time-based triggers that invoke channels at fixed
  intervals. Use these to run periodic actions such as generating a
  stack trace every 5 minutes.

- Channels: Define the action to take and how. For example, the above
  `startup-email` channel defines how to send a message to a specific
  SMTP server.

Channels, signals, and schedules are defined in your `orb` config file.
Each signal maps to a single channel. Multiple signals and schedules can
reference the same channel.

## Channel Details

All channel definitions must have a `name` and a `type`.  Signals
reference channels by `name`.  The channel's `type` must be one of
`notify`, `kafka`, `exec`, `suppress`, `restart`, `kill`, or `signal_toggle`. These
types are described below.

### Channel Configuration Reference

Common fields (may vary by type):
- `name` (string): referenced by signals.
- `type` (string): `notify`, `kafka`, `exec`, `suppress`, `restart`, `kill`, `signal_toggle`.
- `url` (string, notify): destination URL (Go text/template supported).
- `template` (string, notify): message template (optional).
- `broker` (string, kafka): Kafka bootstrap.
- `topic` (string, kafka): topic name.
- `shell` (string, exec): shell script to run.
- `sasl_mechanism` (string, kafka): e.g. `plain` (optional).
- `sasl_username` / `sasl_password` (string, kafka): SASL creds (optional).
- `tls` (bool, kafka): enable TLS.
- `tls_insecure_skip_verify` (bool, kafka): skip verification (not recommended).
- `tls_ca_file` / `tls_cert_file` / `tls_key_file` (string, kafka): TLS files.
- `rate_per_sec` (float): per-channel average actions per second (optional).
- `burst` (int): tokens allowed for bursts (optional; default `1`).

### Sending notifications to messaging platforms

The channel type `notify` is for sending messages to popular messaging
platforms.

You must specify a URL and, optionally, a message template.

`orb` uses `shoutrrr` for sending notifications.  Use the following
URL formats for these different services.  Additional details are
available from the [`shoutrrr` fork documentation](https://github.com/nicholas-fedor/shoutrrr/tree/v0.9.1/docs).

| Service     | URL Format                                                                                 |
|-------------|-------------------------------------------------------------------------------------------- |
| Bark        | `bark://devicekey@host`                                                                    |
| Discord     | `discord://token@id`                                                                       |
| Email       | `smtp://username:password@host:port/?from=fromAddress&to=recipient1[,recipient2,...]`     |
| Gotify      | `gotify://gotify-host/token`                                                               |
| Google Chat | `googlechat://chat.googleapis.com/v1/spaces/FOO/messages?key=bar&token=baz`               |
| IFTTT       | `ifttt://key/?events=event1[,event2,...]&value1=value1&value2=value2&value3=value3`     |
| Join        | `join://shoutrrr:api-key@join/?devices=device1[,device2, ...][&icon=icon][&title=title]` |
| Mattermost  | `mattermost://[username@]mattermost-host/token[/channel]`                                  |
| Matrix      | `matrix://username:password@host:port/[?rooms=!roomID1[,roomAlias2]]`                      |
| Ntfy        | `ntfy://username:password@ntfy.sh/topic`                                                   |
| OpsGenie    | `opsgenie://host/token?responders=responder1[,responder2]`                                 |
| Pushbullet  | `pushbullet://api-token[/device/#channel/email]`                                           |
| Pushover    | `pushover://shoutrrr:apiToken@userKey/?devices=device1[,device2, ...]`                     |
| Rocketchat  | `rocketchat://[username@]rocketchat-host/token[/channel&#124;@recipient]`                  |
| Slack       | `slack://[botname@]token-a/token-b/token-c`                                                |
| Teams       | `teams://group@tenant/altId/groupOwner?host=organization.webhook.office.com`               |
| Telegram    | `telegram://token@telegram?chats=@channel-1[,chat-id-1,...]`                               |
| Zulip Chat  | `zulip://bot-mail:bot-key@zulip-domain/?stream=name-or-id&topic=name`                     |

URLs are actually [go templates](https://pkg.go.dev/text/template)
that are processed before use.  Data that you can access from the
template engine includes:
- `.Timestamp` : the RFC3339 formatted timestamp for the matching log entry
- `.PID`: the process ID for the observed process
- `.Matches`: An array where the first element (`{{index .Matches 0}}`) is the entire log line matched by the regular expression, and subsequent elements (`{{index .Matches 1}}`, `{{index .Matches 2}}`, ...) contain the data captured by each of the regular expression's capturing groups.
- `.Logline`: the matching log line. Equivalent to `{{index .Matches 0}}`.
- `.Env`: a map to access environment variables (e.g. `{{.Env.USER_PASSWORD}}`)

Similarly, the message sent may also be a template.  If no `template`
is specified in the channel definition, then the logline is used as
the message.  If a `template` is specified, then the template is
processed with the same data as above before sending.

As an example, here's a channel that sends an email containing json
data with the observed PID in the email subject line:

```
  - name: email-on-startup
    type: notify
    url:  "smtp://EMAIL@gmail.com:PASSWORD@smtp.gmail.com:587/?from=EMAIL@gmail.com&to=EMAIL@gmail.com&subject=Starting%20process%20{{.PID}}!"
    template: "{ \"timestamp\": \"{{.Timestamp}}\", \"message\": \"{{.Logline}}\" }"
```

Here's an example using environment variables from `.env` for secure configuration:

```
  - name: slack-alerts
    type: notify
    url:  "slack://{{.Env.SLACK_BOT_TOKEN}}@{{.Env.SLACK_CHANNEL}}"
    template: "🚨 Alert from {{.Env.APP_NAME}}: {{.Logline}}"
```

With a `.env` file containing:
```
SLACK_BOT_TOKEN=xoxb-your-bot-token
SLACK_CHANNEL=C1234567890
APP_NAME=my-web-service
```

Generic webhooks are handled specially by `shoutrrr`. Their URL
format is described at
[https://github.com/nicholas-fedor/shoutrrr/blob/v0.9.1/docs/services/generic/index.md](https://github.com/nicholas-fedor/shoutrrr/blob/v0.9.1/docs/services/generic/index.md).

### Sending Kafka messages

The channel type `kafka` is for sending messages to a kafka broker.

You must specify a broker and topic in your definition, like so:

```
  - name: kafka-alerts
    type: kafka
    broker: mybroker.example.com:9092
    topic: orb-alerts
```

Producer timeouts are currently fixed at 5 seconds.

Optional authentication and TLS:

```
  - name: kafka-secure
    type: kafka
    broker: mybroker.example.com:9093
    topic: orb-alerts
    # SASL/PLAIN example
    sasl_mechanism: "plain"
    sasl_username: myuser
    sasl_password: "mypassword"
    # TLS options
    tls: true
    tls_insecure_skip_verify: false  # not recommended in production
    tls_ca_file: "/path/to/ca.pem"   # optional
    tls_cert_file: "/path/to/client.crt"  # optional
    tls_key_file: "/path/to/client.key"   # optional
```

### Running shell scripts

The channel type `exec` is for running arbitrary shell commands.

Environment variables provided to the shell:
- `ORB_PID`: PID of the observed process
- `ORB_MATCH_COUNT`: number of regex matches captured
- `ORB_MATCH_0 .. ORB_MATCH_n`: match values; index 0 is the full line, 1..n are capture groups

In this example, the channel `thread-dump` invokes the `jstack` tool to
dump java thread stacks to a file that is copied into an s3 bucket for
later examination.

```
  - name: thread-dump
    type: exec
    shell: |
      FILENAME=thread-dump-$(date).txt
      jstack "$ORB_PID" > /tmp/${FILENAME}
      aws s3 mv /tmp/${FILENAME} s3:/my-bucket/${FILENAME}

```

You can also reference capture groups from your regex using `ORB_MATCH_1`,
`ORB_MATCH_2`, etc. For example, to echo the first capture group:

```
  - name: dump-first-match
    type: exec
    shell: |
      echo "First match: $ORB_MATCH_1" >> /tmp/matches.txt
```

Compatibility note: older examples used a bash array-like variable to expose
regex matches. This has been replaced by explicit `ORB_MATCH_n` environment
variables for portability and correctness.

## Metrics

Green Orb can expose Prometheus metrics for monitoring throughput and behavior.

- Enable via `--metrics-enable` and set the listen address with `--metrics-addr` (defaults to `127.0.0.1:9090`).
- Exposes `/metrics` with counters, histograms and gauges, including:
  - `orb_events_total{stream}`: lines processed per stream (`stdout|stderr`).
  - `orb_signals_matched_total{signal,channel}`: regex matches.
  - `orb_schedules_fired_total{signal,channel,kind}`: schedule signal firings; kind is `every` or `cron`.
  - `orb_actions_total{channel,type,outcome}`: actions executed and result.
  - `orb_action_latency_seconds{channel,type}`: action latency.
  - `orb_dropped_events_total{reason}`: dropped events (`queue_full|rate_limited`).
  - `orb_queue_depth`: current queue size.
  - `orb_observed_pid`: PID of the observed process.

Example:

```
orb --metrics-enable --metrics-addr 127.0.0.1:9090 myapp ...
```

With environment file:

```
orb --env production.env --metrics-enable java -jar myapp.jar
```

Prometheus scrape example:

```
scrape_configs:
  - job_name: 'green-orb'
    static_configs:
      - targets: ['127.0.0.1:9090']
```

## Rate Limiting and Non-Blocking Queue

To prevent backpressure and notification storms, Green Orb uses a non-blocking queue and supports per-channel rate limiting.

- Non-blocking enqueue drops when the queue is full and counts drops in `orb_dropped_events_total{reason="queue_full"}`.
- Configure per-channel token-bucket rate limits with optional `rate_per_sec` and `burst` fields:

```
channels:
  - name: slack-alerts
    type: notify
    url:  "slack://..."
    rate_per_sec: 2.0   # average 2 actions per second
    burst: 5            # allow short bursts
```

Events that exceed rate limits are dropped and counted in `orb_dropped_events_total{reason="rate_limited"}`.

Operational tips:
- If you see `queue_full` drops, increase `--workers`, add per-channel `rate_per_sec`, or reduce event volume.
- Keep `/metrics` bound to localhost or protect it behind auth when needed.

## Checks (HTTP/TCP/Flapping)

In addition to log-driven signals, Green Orb can run periodic checks and trigger channels when conditions fail.

- Define checks in `checks:`; each check sends a message to a `channel` when it fails.
- Supported types: `http`, `tcp`, `flapping` (restart burst detection).

Examples:

```
checks:
  - name: web-status
    type: http
    url: https://example.com/health
    expect_status: 200
    body_regex: OK           # optional
    interval: 30s
    timeout: 5s
    channel: email_alerts

  - name: db-port
    type: tcp
    host: db.internal
    port: 5432
    interval: 15s
    timeout: 3s
    channel: email_alerts

  - name: service-flapping
    type: flapping
    restart_threshold: 3        # restarts
    window: 5m                  # within this window
    interval: 30s
    channel: email_alerts
```

Metrics: `orb_checks_total{type,outcome}` increments on each run with `outcome`=`success|error`.

### Suppressing output

The channel type `suppress` is for suppressing output from your
observed process.  Anything sent to a `suppress` channel will not flow
through to standard output or standard error.

### Restarting your process

The channel type `restart` is for restarting your observed process.

The `orb` process will run continuously, but it will force the
observed process to terminate and then restart.

### Killing your process

The channel type `kill` is for killing your observed process.

The `orb` process will exit.

## Contributing

Green Orb is Free Software, and any and all contributions are welcome!
Please use GitHub's Issue tracker and Pull Request systems for
feedback and improvements.

## Author and License

Green Orb was written by [Anthony
Green](https://github.com/atgreen), and is distributed under the terms
of the MIT License.  See
[LICENSE](https://raw.githubusercontent.com/atgreen/green-orb/main/LICENSE)
for details.
