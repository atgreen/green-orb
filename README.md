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
  - name: "slack_errors"
    type: "slack"
    settings:
      token: "xoxb-your-slack-token"
      channel: "#errors"

signals:
  - regex: "^ERROR:"
    channel: "slack_errors"
```

`oarb` does not interfere with the execution of your program.  All
console logs still go to the console, and the exit code for your
program is passed on through `oarb`.

Author and License
-------------------

`oarb` was written by [Anthony
Green](https://github.com/atgreen), and is distributed under the terms
of the MIT License.  See
[LICENSE](https://raw.githubusercontent.com/atgreen/oarb/main/LICENSE)
for details.
