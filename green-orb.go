// green-orb.go - an Observe and Report Buddy
//
// SPDX-License-Identifier: MIT
//
// Copyright (C) 2023-2025 Anthony Green - green@moxielogic.com
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use, copy,
// modify, merge, publish, distribute, sublicense, and/or sell copies
// of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
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
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"

	"github.com/urfave/cli/v3"
)

var version = "dev"

func main() {
	var configFilePath string
	var numWorkers int64
	var metricsAddr string

	cmd := &cli.Command{
		Name:            "orb",
		HideHelpCommand: true,
		Version:         version,
		Usage:           "Your observe-and-report buddy",
		Copyright:       "Copyright (C) 2023-2025  Anthony Green <green@moxielogic.com>.\nDistributed under the terms of the MIT license.\nSee https://github.com/atgreen/green-orb for details.",
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
			&cli.BoolFlag{
				Name:        "metrics-enable",
				Value:       false,
				Usage:       "enable Prometheus metrics endpoint",
				Destination: &metricsEnable,
			},
			&cli.StringFlag{
				Name:        "metrics-addr",
				Value:       "127.0.0.1:9090",
				Usage:       "Prometheus metrics listen address",
				Destination: &metricsAddr,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.NArg() == 0 {
				// No arguments provided, show help text
				cli.ShowAppHelp(cmd)
				return nil
			}
			return runObserved(ctx, configFilePath, numWorkers, metricsAddr, cmd.Args().Slice())
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

// runObserved contains the core execution logic for running and observing a subprocess.
func runObserved(ctx context.Context, configFilePath string, numWorkers int64, metricsAddr string, subprocessArgs []string) error {
	// Drop any leading "--" separators that may be present after CLI parsing.
	for len(subprocessArgs) > 0 && subprocessArgs[0] == "--" {
		subprocessArgs = subprocessArgs[1:]
	}

	// Load and validate configuration
	config, err := LoadConfig(configFilePath)
	if err != nil {
		log.Fatalf("green-orb error: Failed to load config: %v", err)
	}

	// Initialize Kafka manager
	kafkaManager := NewKafkaManager()
	if err := kafkaManager.Connect(config.Channels); err != nil {
		log.Fatalf("green-orb error: Failed to initialize Kafka: %v", err)
	}
	defer kafkaManager.Close()

	// Compile signals
	compiledSignals, err := CompileSignals(config.Signals)
	if err != nil {
		log.Fatalf("green-orb error: Failed to compile signals: %v", err)
	}

	if len(subprocessArgs) == 0 {
		log.Fatal("green-orb error: No command provided to run")
	}

	// Initialize metrics if enabled
	if metricsEnable {
		StartMetricsServer(metricsAddr)
	}

	// Create channel map
	channelMap := CreateChannelMap(config.Channels)
	channels = channelMap

	// Create worker pool
	workerPool := NewWorkerPool(int(numWorkers), 100, channelMap, kafkaManager.clients)
	workerPool.Start()
	defer workerPool.Stop()

	// Start checks scheduler if any checks are defined
	checkScheduler := NewCheckScheduler(config.Checks, func() int {
		if observedCmd != nil && observedCmd.Process != nil {
			return observedCmd.Process.Pid
		}
		return 0
	}, channelMap, workerPool)
	if len(config.Checks) > 0 {
		checkScheduler.Start()
		defer checkScheduler.Stop()
	}

	// Process restart loop
	for shouldRestart {
		restartMutex.Lock()
		shouldRestart = false
		restartMutex.Unlock()

		// Prepare to run the subprocess
		observedCmd = exec.Command(subprocessArgs[0], subprocessArgs[1:]...)
		stdout, err := observedCmd.StdoutPipe()
		if err != nil {
			log.Fatalf("green-orb error: Failed to create stdout pipe: %v", err)
		}
		stderr, err := observedCmd.StderrPipe()
		if err != nil {
			log.Fatalf("green-orb error: Failed to create stderr pipe: %v", err)
		}

		// Set up signal handling
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan)

		if err := observedCmd.Start(); err != nil {
			log.Fatalf("green-orb error: Failed to start subprocess: %v", err)
		}

		// Goroutine for passing signals
		go func() {
			for sig := range sigChan {
				process_signal(observedCmd, sig)
			}
		}()

		pid := observedCmd.Process.Pid
		if metricsEnable {
			orbObservedPID.Set(float64(pid))
		}

		// Create monitor and start monitoring
		monitor := NewMonitor(pid, compiledSignals, workerPool, channelMap)

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			monitor.MonitorOutput(bufio.NewScanner(stdout), false)
		}()
		go func() {
			defer wg.Done()
			monitor.MonitorOutput(bufio.NewScanner(stderr), true)
		}()

		// Wait for all reading to be complete
		wg.Wait()

		// Wait for the command to finish
		err = observedCmd.Wait()

		// After cmd.Wait(), stop listening for signals
		signal.Stop(sigChan)
		close(sigChan)

		// Check if we should restart
		restartMutex.Lock()
		restart := shouldRestart
		restartMutex.Unlock()

		if !restart {
			break
		}
	}

	if metricsEnable {
		orbObservedPID.Set(0)
	}

	// Handle exit status
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			os.Exit(exitError.ExitCode())
		}
		log.Fatalf("green-orb error: Error waiting for Command: %v", err)
	}

	return nil
}