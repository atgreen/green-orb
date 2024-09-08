// green-orb.go - an Observe and Report Buddy
//
// SPDX-License-Identifier: MIT
//
// Copyright (C) 2023, 2024 Anthony Green - green@moxielogic.com
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the “Software”), to deal in the Software without
// restriction, including without limitation the rights to use, copy,
// modify, merge, publish, distribute, sublicense, and/or sell copies
// of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
// BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
// ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/containrrr/shoutrrr"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/urfave/cli/v3"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strings"
	"sync"
	"syscall"
	"text/template"
	"time"
)

var version = "dev"
var restart bool = true
var observed_cmd *exec.Cmd

type Channel struct {
	Name     string `yaml:"name"`
	Type     string `yaml:"type"`
	URL      string `yaml:"url"`
	Template string `yaml:"template"`
	Topic    string `yaml:"topic"`
	Broker   string `yaml:"broker"`
	Shell    string `yaml:"shell"`
}

type Signal struct {
	Regex   string `yaml:"regex"`
	Channel string `yaml:"channel"`
}

type Config struct {
	Channel []Channel `yaml:"channels"`
	Signal  []Signal  `yaml:"signals"`
}

var kafkaClients map[string]*kgo.Client = make(map[string]*kgo.Client)

func kafkaConnect(channels []Channel) (map[string]*kgo.Client, error) {

	for _, channel := range channels {
		if channel.Type == "kafka" {
			seeds := []string{channel.Broker}
			opts := []kgo.Opt{
				kgo.RequiredAcks(kgo.AllISRAcks()),
				kgo.DisableIdempotentWrite(),
				kgo.ProducerLinger(50 * time.Millisecond),
				kgo.RecordRetries(math.MaxInt32),
				kgo.RecordDeliveryTimeout(5 * time.Second),
				kgo.ProduceRequestTimeout(5 * time.Second),
				kgo.SeedBrokers(seeds...),
			}

			cl, err := kgo.NewClient(opts...)
			if err != nil {
				log.Fatal("green-orb error: failed to create kafka client connection: ", err)
			}
			kafkaClients[channel.Name] = cl
		}
	}
	return kafkaClients, nil
}

func compileSignals(signals []Signal) ([]CompiledSignal, error) {
	var compiledSignals []CompiledSignal
	for _, signal := range signals {
		re, err := regexp.Compile(signal.Regex)
		if err != nil {
			log.Fatal("green-orb error: failed to compile regex \"", signal.Regex, "\": ", err)
		}
		compiledSignals = append(compiledSignals, CompiledSignal{
			Regex:   re,
			Channel: signal.Channel,
		})
	}
	return compiledSignals, nil
}

type CompiledSignal struct {
	Regex   *regexp.Regexp
	Channel string
}

type Notification struct {
	PID     int
	Channel Channel
	Match  []string
	Message string
}

type TemplateData struct {
	PID       int
	Logline   string
	Matches   []string
	Timestamp string
	Env       map[string]string
}

func loadYAMLConfig(filename string, config *Config) error {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("green-orb error: ", err)
	}

	err = yaml.Unmarshal(bytes, config)
	if err != nil {
		log.Fatal("green-orb error: Failed parsing config file: ", err)
	}

	return nil
}

func envToMap() (map[string]string, error) {
	envMap := make(map[string]string)
	var err error

	for _, v := range os.Environ() {
		split_v := strings.Split(v, "=")
		envMap[split_v[0]] = strings.Join(split_v[1:], "=")
	}
	return envMap, err
}

func startWorkers(notificationQueue <-chan Notification, numWorkers int64, wg *sync.WaitGroup) {
	var messageString string

	for i := 0; i < int(numWorkers); i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for notification := range notificationQueue {
				env, _ := envToMap()
				td := TemplateData{PID: notification.PID,
					Logline:   notification.Message,
					Env:       env,
					Matches:   notification.Match,
					Timestamp: time.Now().Format(time.RFC3339)}
				switch notification.Channel.Type {
				case "notify":
					tmpl, err := template.New("url").Parse(notification.Channel.URL)
					if err != nil {
						log.Fatal("green-orb error: can't parse URL template: ", err)
					}
					var buffer bytes.Buffer
					err = tmpl.Execute(&buffer, td)
					if err != nil {
						log.Fatal("green-orb error: can't execute URL template: ", err)
					}
					urlString := buffer.String()
					if notification.Channel.Template != "" {
						tmpl, err := template.New("msg").Parse(notification.Channel.Template)
						if err != nil {
							log.Fatal("green-orb error: can't parse template: ", err)
						}
						var buffer bytes.Buffer
						err = tmpl.Execute(&buffer, td)
						if err != nil {
							log.Fatal("green-orb error: can't execute URL template: ", err)
						}
						messageString = buffer.String()
					} else {
						messageString = notification.Message
					}
					err = shoutrrr.Send(urlString, messageString)
					if err != nil {
						log.Println("green-orb warning: failed sending notification: ", err)
					}
				case "exec":
					var stdout bytes.Buffer
					var stderr bytes.Buffer

					serializedMatches := "("
					for _, m := range notification.Match[0:] {
						encoded := base64.StdEncoding.EncodeToString([]byte(m))
						serializedMatches += fmt.Sprintf("$(echo %s | base64 -d) ", encoded)
					}
					serializedMatches += ")"

					cmd := exec.Command("bash", "-c", "export ORB_MATCHES=" + serializedMatches + "; " + notification.Channel.Shell)
					cmd.Env = os.Environ()
					cmd.Env = append(cmd.Env, fmt.Sprintf("ORB_PID=%d", notification.PID))
					cmd.Stdout = &stdout
					cmd.Stderr = &stderr
					cmd.Run()
				case "kafka":
					ctx := context.Background()
					record := &kgo.Record{Topic: notification.Channel.Topic, Value: []byte(notification.Message)}
					if err := kafkaClients[notification.Channel.Name].ProduceSync(ctx, record).FirstErr(); err != nil {
						log.Println("green-orb warning: kafka record had a produce error:", err)
					}
				case "restart":
					restart = true
					observed_cmd.Process.Signal(syscall.SIGTERM)
				case "kill":
					restart = false
					observed_cmd.Process.Signal(syscall.SIGTERM)
				}
			}
		}(i)
	}
}

func main() {

	var configFilePath string
	var numWorkers int64

	cmd := &cli.Command{
		Name:            "orb",
		HideHelpCommand: true,
		Version:         version,
		Usage:           "Your observe-and-report buddy",
		Copyright:       "Copyright (C) 2023-2024  Anthony Green <green@moxielogic.com>.\nDistributed under the terms of the MIT license.\nSee https://github.com/atgreen/green-orb for details.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Value:       "green-orb.yaml",
				Aliases:     []string{"c"},
				Usage:       "path to the green-orb configuration file",
				Destination: &configFilePath,
			},
			&cli.IntFlag{
				Name:    "workers",
				Value:   5,
				Aliases: []string{"w"},
				Usage:   "number of reporting workers",
				Action: func(ctx context.Context, cmd *cli.Command, v int64) error {
					if (v > 100) || (v < 1) {
						return fmt.Errorf("Flag workers value %v out of range [1-100]", v)
					}
					return nil
				},
				Destination: &numWorkers,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {

			if cmd.NArg() == 0 {
				// No arguments provided, show help text
				cli.ShowAppHelp(cmd)
				return nil
			}

			config := Config{}
			err := loadYAMLConfig(configFilePath, &config)
			if err != nil {
				log.Fatal("green-orb error: Failed to load config: ", err)
			}

			kafkaConnect(config.Channel)
			compiledSignals, _ := compileSignals(config.Signal)

			// The remaining arguments after flags are parsed
			subprocessArgs := cmd.Args().Slice()
			if len(subprocessArgs) == 0 {
				log.Fatal("green-orb error: No command provided to run")
			}

			notificationQueue := make(chan Notification, 100)

			// Use a WaitGroup to wait for the reading goroutines to finish
			var wg sync.WaitGroup
			var nwg sync.WaitGroup

			startWorkers(notificationQueue, numWorkers, &nwg)

			channelMap := make(map[string]Channel)
			for _, ch := range config.Channel {
				channelMap[ch.Name] = ch
			}

			for restart {

				restart = false

				// Prepare to run the subprocess
				observed_cmd = exec.Command(subprocessArgs[0], subprocessArgs[1:]...)
				// Rest of your code to handle the subprocess...
				stdout, _ := observed_cmd.StdoutPipe()
				stderr, _ := observed_cmd.StderrPipe()

				sigChan := make(chan os.Signal, 1)
				signal.Notify(sigChan)

				if err := observed_cmd.Start(); err != nil {
					log.Fatal("green-orb error: Failed to start subprocess: ", err)
				}

				// Goroutine for passing signals
				go func() {
					for sig := range sigChan {
						process_signal(observed_cmd, sig)
					}
				}()

				pid := observed_cmd.Process.Pid

				// Increment WaitGroup and start reading stdout
				wg.Add(2)
				go func() {
					defer wg.Done()
					monitorOutput(pid, bufio.NewScanner(stdout), compiledSignals, notificationQueue, channelMap, false)
				}()
				go func() {
					defer wg.Done()
					monitorOutput(pid, bufio.NewScanner(stderr), compiledSignals, notificationQueue, channelMap, true)
				}()

				// Wait for all reading to be complete
				wg.Wait()

				// Wait for the command to finish
				err = observed_cmd.Wait()

				// After cmd.Wait(), stop listening for signals
				signal.Stop(sigChan)
				close(sigChan)
			}

			close(notificationQueue)

			// Wait for all reading to be complete
			nwg.Wait()

			// Handle exit status
			if err != nil {
				// observed_cmd.Wait() returns an error if the command exits non-zero
				if exitError, ok := err.(*exec.ExitError); ok {
					// Get the command's exit code
					os.Exit(exitError.ExitCode())
				} else {
					// Other error types (not non-zero exit)
					log.Fatal("green-orb error: Error waiting for Command:", err)
				}
			} else {
				// Success (exit code 0)
				os.Exit(0)
			}
			return nil
		},
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func monitorOutput(pid int, scanner *bufio.Scanner, compiledSignals []CompiledSignal, notificationQueue chan Notification, channelMap map[string]Channel, is_stderr bool) {
	for scanner.Scan() {
		line := scanner.Text()
		suppress := false;

		for _, signal := range compiledSignals {
			match := signal.Regex.FindStringSubmatch(line)
			if (match != nil) {
				channel, _ := channelMap[signal.Channel]
				if channel.Type == "suppress" {
					suppress = true
				} else {
					notificationQueue <- Notification{PID: pid, Match: match, Channel: channel, Message: line}
				}
			}
		}

	  if (! suppress) {
			if is_stderr {
				fmt.Fprintln(os.Stderr, line)
			} else {
				fmt.Println(line)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("green-orb error: Problem reading from pipe: ", err)
	}
}
