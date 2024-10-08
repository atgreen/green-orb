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

This assumes `green-orb.yaml` is in the current directory.  Use the
`-c` flag to point it elsewhere.

The config file tells `orb` what to watch for, and what to do.  For
instance, if it contains the following, you'll get an email every time
your application starts up, and a thread dump every time you get a
thread pool exhausted warning.

```
channels:
  - name: "startup-email"
    type: "notify"
    url:  "smtp://MYEMAIL@gmail.com:MYPASSWORD@smtp.gmail.com:587/?from=MYEMAIL@gmail.com&to=MYEMAIL@gmail.com&subject=Application%20Starting!"

  - name: "thread-dump"
    type: "exec"
    shell: |
      FILENAME=thread-dump-$(date).txt
      jstack $ORB_PID > /tmp/${FILENAME}
      aws s3 mv /tmp/${FILENAME} s3:/my-bucket/${FILENAME}

signals:
  - regex: "Starting Application"
    channel: "startup-email"

  - regex: "Warning: thread pool exhausted"
    channel: "thread-dump"
```

`orb` does not interfere with the execution of your application.  All
console logs still go to the console, Linux and macOS signals
are passed through to the observed process, and the exit code for your
application is returned through `orb`.

## Channels and Signals

In Green Orb, "channels" and "signals" are foundational concepts:

- Signals: Mappings of regular expressions to channels. When a
  signal's regex matches a log entry, the corresponding channel is
  invoked.

- Channels: Define the action to take and how. For example, the above
  `startup-email` channel defines how to send a message to a specific
  SMTP server.

Channels and signals are defined in your `orb` config file.  A config
file can define any number of channels and signals.  Each signal maps
to a single channel.  However, multiple signals can map to the same
channel.

## Channel Details

All channel definitions must have a `name` and a `type`.  Signals
reference channels by `name`.  The channel's `type` must be one of
`notify`, `kafka`, `exec`, `suppress`, `restart` or `kill`.  These
types are described below.

### Sending notifications to messaging platforms

The channel type `notify` is for sending messages to popular messaging
platforms.

You must specify a URL and, optionally, a message template.

`orb` uses `shoutrrr` for sending notifications.  Use the following
URL formats for these different services.  Additional details are
available from the [`shoutrrr`
documentation](https://containrrr.dev/shoutrrr/v0.8/services/overview/).

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
  - name: "email-on-startup"
    type: "notify"
    url:  "smtp://EMAIL@gmail.com:PASSWORD@smtp.gmail.com:587/?from=EMAIL@gmail.com&to=EMAIL@gmail.com&subject=Starting%20process%20{{.PID}}!"
    template: "{ \"timestamp\": \"{{.Timestamp}}\", \"message\": \"{{.Logline}}\" }"
```

Generic webhooks and handled specially by `shoutrrr`.  Their URL
format is described at
[https://containrrr.dev/shoutrrr/v0.8/services/generic/](https://containrrr.dev/shoutrrr/v0.8/services/generic/).

### Sending Kafka messages

The channel type `kafka` is for sending messages to a kafka broker.

You must specify a broker and topic in your definition, like so:

```
  - name: "kafka-alerts"
    type: "kafka"
    broker: "mybroker.example.com:9092"
    topic: "orb-alerts"
```

Producer timeouts are currently fixed at 5 seconds.

### Running shell scripts

The channel type `exec` is for running arbitrary shell commands.

The process ID of the observed process is presented to the shell code
through the environment variable `$ORB_PID`.  In this example, the
channel `thread-dump` invokes the `jstack` tool to dump java thread
stacks to a file that is copied into an s3 bucket for later
examination.

```
  - name: "thread-dump"
    type: "exec"
    shell: |
      FILENAME=thread-dump-$(date).txt
      jstack $ORB_PID > /tmp/${FILENAME}
      aws s3 mv /tmp/${FILENAME} s3:/my-bucket/${FILENAME}
```

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
