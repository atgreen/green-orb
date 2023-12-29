# oarb
> An 'Observe and Report Buddy' for your SRE toolbox

`oarb` monitors your program's console output for regexps that you
define, and sends notifications to specific channels on every match.

Supported notification channels include:
- webhooks
- slack

`oarb` is very easy to configure and use.  It is one binary and one yaml
config file.  Simply add `oarb -c config.yaml` before your program.  For instance, instead of:
```
$ java -jar mywebapp.jar
```
...use...
```
$ oarb -c config.yaml java -jar mywebapp.jar
```

If `config.yaml` contains the following, you'll get a slack message on
every console log message that starts with `ERROR:`:

```
channels:
  - name: "slack_alerts"
    type: "sender"
    url:  "generic:localhost:5000/?contentType=json"

signals:
  - regex: "^ERROR:"
    channel: "slack_alerts"
```

`oarb` does not interfere with the execution of your program.  All
console logs still go to the console, and the exit code for your
program is passed on through `oarb`.

Sender Services
----------------

`oarb` uses `shoutrrr` for sending notifications.  Use the following
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


Author and License
-------------------

`oarb` was written by [Anthony
Green](https://github.com/atgreen), and is distributed under the terms
of the MIT License.  See
[LICENSE](https://raw.githubusercontent.com/atgreen/oarb/main/LICENSE)
for details.
