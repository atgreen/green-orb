// orb-ag.go - an Observe and Report Buddy
//
// Copyright (C) 2023 Anthony Green - green@moxielogic.com
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
	"regexp"
	"sync"
	"syscall"
	"time"
)

var version = "dev"
var restart bool = true
var observed_cmd *exec.Cmd

type Channel struct {
	Name   string `yaml:"name"`
	Type   string `yaml:"type"`
	URL    string `yaml:"url"`
	Topic  string `yaml:"topic"`
	Broker string `yaml:"broker"`
	Shell  string `yaml:"shell"`
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
				log.Fatal("orb-ag error: failed to create kafka client connection: ", err)
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
			log.Fatal("orb-ag error: failed to compile regex \"", signal.Regex, "\": ", err)
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
	Message string
}

func loadYAMLConfig(filename string, config *Config) error {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("orb-ag error: ", err)
	}

	err = yaml.Unmarshal(bytes, config)
	if err != nil {
		log.Fatal("orb-ag error: Failed parsing config file: ", err)
	}

	return nil
}

func startWorkers(notificationQueue <-chan Notification, numWorkers int, wg *sync.WaitGroup) {
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for notification := range notificationQueue {
				switch notification.Channel.Type {
				case "notify":
					err := shoutrrr.Send(notification.Channel.URL, notification.Message)
					if err != nil {
						log.Println("org-ab warning: failed sending notification: ", err)
					}
				case "exec":
					var stdout bytes.Buffer
					var stderr bytes.Buffer
					cmd := exec.Command("bash", "-c", notification.Channel.Shell)
					cmd.Env = os.Environ()
					cmd.Env = append(cmd.Env, fmt.Sprintf("ORB_PID=%d", notification.PID))
					cmd.Stdout = &stdout
					cmd.Stderr = &stderr
					cmd.Run()
				case "kafka":
					ctx := context.Background()
					record := &kgo.Record{Topic: notification.Channel.Topic, Value: []byte(notification.Message)}
					if err := kafkaClients[notification.Channel.Name].ProduceSync(ctx, record).FirstErr(); err != nil {
						log.Println("orb-ag warning: kafka record had a produce error:", err)
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

	cmd := &cli.Command{
		Name:            "orb-ag",
		HideHelpCommand: true,
		Version:         version,
		Usage:           "Your observe-and-report buddy",
		Copyright:       "Copyright (C) 2023  Anthony Green <green@moxielogic.com>.\nDistributed under the terms of the MIT license.\nSee https://github.com/atgreen/orb-ag for details.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Value:       "orb-ag.yaml",
				Aliases:     []string{"c"},
				Usage:       "path to the orb-ag configuration file",
				Destination: &configFilePath,
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
				log.Fatal("orb-ag error: Failed to load config: ", err)
			}

			kafkaConnect(config.Channel)
			compiledSignals, _ := compileSignals(config.Signal)

			// The remaining arguments after flags are parsed
			subprocessArgs := cmd.Args().Slice()
			if len(subprocessArgs) == 0 {
				log.Fatal("orb-ag error: No command provided to run")
			}

			notificationQueue := make(chan Notification, 100)

			// Use a WaitGroup to wait for the reading goroutines to finish
			var wg sync.WaitGroup
			var nwg sync.WaitGroup

			startWorkers(notificationQueue, 5, &nwg)

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

				if err := observed_cmd.Start(); err != nil {
					// Handle error
				}

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
					log.Fatal("orb-ag error: Error waiting for Command:", err)
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

		for _, signal := range compiledSignals {
			if signal.Regex.MatchString(line) {
				channel, _ := channelMap[signal.Channel]
				notificationQueue <- Notification{PID: pid, Channel: channel, Message: line}
			}
		}

		if is_stderr {
			fmt.Fprintln(os.Stderr, line)
		} else {
			fmt.Println(line)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("orb-ag error: Problem reading from pipe: ", err)
	}
}
