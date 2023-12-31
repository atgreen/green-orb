// orb-ag.go - Observe and Report Buddy
//
// Copyright (C) 2023 Anthony Green - green@moxielogic.com
//

package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"github.com/containrrr/shoutrrr"
	"github.com/twmb/franz-go/pkg/kgo"
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
		log.Fatal("orb-ag error: Failed reading config file: ", err)
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
						log.Println("org-ab warning: failed sending notification: %s\n", err)
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

	// Define a string flag for the configuration file
	configFilePath := flag.String("c", "", "Path to the oar configuration file")

	// Parse the flags
	flag.Parse()

	config := Config{}
	err := loadYAMLConfig(*configFilePath, &config)
	if err != nil {
		log.Fatal("orb-ag error: Failed to load config: ", err)
	}

	kafkaConnect(config.Channel)
	compiledSignals, _ := compileSignals(config.Signal)

	// The remaining arguments after flags are parsed
	subprocessArgs := flag.Args()
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
			monitorOutput(pid, bufio.NewScanner(stdout), compiledSignals, notificationQueue, channelMap)
		}()
		go func() {
			defer wg.Done()
			monitorOutput(pid, bufio.NewScanner(stderr), compiledSignals, notificationQueue, channelMap)
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
		// cmd.Wait() returns an error if the command exits non-zero
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
}

func monitorOutput(pid int, scanner *bufio.Scanner, compiledSignals []CompiledSignal, notificationQueue chan Notification, channelMap map[string]Channel) {
	for scanner.Scan() {
		line := scanner.Text()

		for _, signal := range compiledSignals {
			if signal.Regex.MatchString(line) {
				channel, _ := channelMap[signal.Channel]
				notificationQueue <- Notification{PID: pid, Channel: channel, Message: line}
			}
		}

		fmt.Println(line)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("orb-ag error: Problem reading from pipe: %v\n", err)
	}
}
