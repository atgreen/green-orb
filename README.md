# orb-ag
> An 'Observe and Report Buddy' for your SRE toolbox

`orb-ag` monitors your program's console output for patterns that you
define, and performs actions based on what it detects.  Actions
include:

* sending notifications through one of a number of popular messaging
  services, including slack, email, discord, and many more.
* sending webhooks for processing by services like the Ansible
  Automation Platform.
* sending messages on a kafka topic
* executing arbitrary shell commands, allowing you to, for instance,
  capture thread dumps of the process being observed.
* restarting the program being observed (avoiding pod restarts on k8s
  platforms, for instance)

`orb-ag` is very easy to configure and use.  It's just one binary and one yaml
config file.  Simply preface your program with `orb-ag -c config.yaml`.  For example, instead of:
```
$ java -jar mywebapp.jar
```
...use...
```
$ orb-ag -c config.yaml java -jar mywebapp.jar
```

Or, if you are using containers, change...
```
ENTRYPOINT [ "java","-jar", "jar-file-name.jar" ]
```
...in your Dockerfile, to...
```
ENTRYPOINT [ "orb-ag", "-c", "config.yaml", "java","-jar", "jar-file-name.jar" ]
```



If `config.yaml` contains the following, you'll get a slack message on
every console log message that starts with `ERROR:`:

```
channels:
  - name: "slack_alerts"
    type: "notify"
    url:  "generic:localhost:5000/?contentType=json"

signals:
  - regex: "^ERROR:"
    channel: "slack_alerts"
```

`orb-ag` does not interfere with the execution of your program.  All
console logs still go to the console, and the exit code for your
program is passed on through `orb-ag`.

# Action Details

## Sending notifications to messaging platforms

The channel type `notify` is for sending messages to popular messaging
platforms.

`orb-ag` uses `shoutrrr` for sending notifications.  Use the following
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

The URL format for generic webhooks is described at
[https://containrrr.dev/shoutrrr/v0.8/services/generic/](https://containrrr.dev/shoutrrr/v0.8/services/generic/).

## Sending Kafka messages

The channel type `kafka` is for sending messages to a kafka broker.

## Running shell scripts

The channel type `exec` is for running arbitrary shell commands.

The process ID of the observed process is presented to the shell code
through the environment variable `$ORB_PID`.  In this example, the
channel `thread-dump` invokes the `jstack` tool to dump java thread
stacks to a temporary file for later examination.

```
  - name: "thread-dump"
    type: "exec"
    shell: |
      jstack $ORB_PID > /tmp/thread-dump-$(date).txt 2>&1
```

## Restarting your process

Author and License
-------------------

`orb-ag` was written by [Anthony
Green](https://github.com/atgreen), and is distributed under the terms
of the MIT License.  See
[LICENSE](https://raw.githubusercontent.com/atgreen/orb-ag/main/LICENSE)
for details.
