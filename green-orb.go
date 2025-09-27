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
    "errors"
    "fmt"
    "log"
    "os"
    "os/exec"
    "os/signal"
    "strings"
    "sync"
)

var version = "dev"

// loadEnvFile loads environment variables from a .env file
func loadEnvFile(filepath string) error {
    file, err := os.Open(filepath)
    if err != nil {
        return err
    }
    defer func() {
        _ = file.Close()
    }()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse KEY=VALUE format
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue // Skip malformed lines
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes if present
		if len(value) >= 2 {
			if (strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`)) ||
				(strings.HasPrefix(value, `'`) && strings.HasSuffix(value, `'`)) {
				value = value[1 : len(value)-1]
			}
		}

        // Set environment variable
        if err := os.Setenv(key, value); err != nil {
            log.Printf("green-orb warning: failed to set env %s: %v", key, err)
        }
	}

	return scanner.Err()
}

func main() {
	var configFilePath = "green-orb.yaml"
	var numWorkers int64 = 5
	var metricsAddr = "127.0.0.1:9090"
	var envFile = ""
	var skipDotEnv = false

	// Parse arguments manually to stop at first non-flag
	args := os.Args[1:] // Skip program name
	var commandArgs []string

	for i := 0; i < len(args); i++ {
		arg := args[i]

		// Check for help/version flags
		if arg == "--help" || arg == "-h" {
			showHelp()
			return
		}
		if arg == "--version" || arg == "-v" {
			fmt.Printf("orb version %s\n", version)
			return
		}

		// Stop parsing at first non-flag
		if !strings.HasPrefix(arg, "-") {
			commandArgs = args[i:]
			break
		}

		// Handle flags
		if arg == "--config" || arg == "-c" {
			if i+1 >= len(args) {
				log.Fatal("green-orb error: --config requires a value")
			}
			configFilePath = args[i+1]
			i++ // Skip the value
		} else if arg == "--workers" || arg == "-w" {
			if i+1 >= len(args) {
				log.Fatal("green-orb error: --workers requires a value")
			}
			var err error
			if numWorkers, err = parseWorkers(args[i+1]); err != nil {
				log.Fatal("green-orb error: ", err)
			}
			i++ // Skip the value
		} else if arg == "--metrics-enable" {
			metricsEnable = true
		} else if arg == "--metrics-addr" {
			if i+1 >= len(args) {
				log.Fatal("green-orb error: --metrics-addr requires a value")
			}
			metricsAddr = args[i+1]
			i++ // Skip the value
		} else if arg == "--env" {
			if i+1 >= len(args) {
				log.Fatal("green-orb error: --env requires a value")
			}
			envFile = args[i+1]
			i++ // Skip the value
		} else if arg == "--skip-dotenv" {
			skipDotEnv = true
		} else if strings.HasPrefix(arg, "--config=") {
			configFilePath = arg[9:] // Remove "--config="
		} else if strings.HasPrefix(arg, "--workers=") {
			var err error
			if numWorkers, err = parseWorkers(arg[10:]); err != nil {
				log.Fatal("green-orb error: ", err)
			}
		} else if strings.HasPrefix(arg, "--metrics-addr=") {
			metricsAddr = arg[15:] // Remove "--metrics-addr="
		} else if strings.HasPrefix(arg, "--env=") {
			envFile = arg[6:] // Remove "--env="
		} else {
			log.Fatalf("green-orb error: unknown flag: %s", arg)
		}
	}

	if len(commandArgs) == 0 {
		showHelp()
		return
	}

    if err := runObserved(configFilePath, numWorkers, metricsAddr, envFile, !skipDotEnv, commandArgs); err != nil {
        log.Fatal(err)
    }
}

func parseWorkers(value string) (int64, error) {
	var w int64
	if _, err := fmt.Sscanf(value, "%d", &w); err != nil {
		return 0, fmt.Errorf("invalid workers value: %s", value)
	}
	if w < 1 || w > 100 {
		return 0, fmt.Errorf("workers value %d out of range [1-100]", w)
	}
	return w, nil
}

func showHelp() {
	fmt.Printf(`NAME:
   orb - Your observe-and-report buddy

USAGE:
   orb [global options] command [command arguments...]

VERSION:
   %s

GLOBAL OPTIONS:
   --config value, -c value   path to the green-orb configuration file (default: "green-orb.yaml")
   --workers value, -w value  number of reporting workers (default: 5)
   --metrics-enable           enable Prometheus metrics endpoint (default: false)
   --metrics-addr value       Prometheus metrics listen address (default: "127.0.0.1:9090")
   --env value                load environment variables from specified file
   --skip-dotenv              do not automatically load .env file (default: loads .env if present)
   --help, -h                 show help
   --version, -v              print the version

EXAMPLES:
   orb echo "Hello World"                     # Observe echo command (loads .env if present)
   orb -c myconfig.yaml ls -la                # Use custom config with ls -la
   orb --metrics-enable java -jar app.jar     # Enable metrics while observing Java app
   orb --skip-dotenv npm start                # Skip loading .env file
   orb --env production.env node app.js       # Load custom env file (plus .env unless --skip-dotenv)

COPYRIGHT:
   Copyright (C) 2023-2025  Anthony Green <green@moxielogic.com>.
   Distributed under the terms of the MIT license.
   See https://github.com/atgreen/green-orb for details.
`, version)
}

// runObserved contains the core execution logic for running and observing a subprocess.
func runObserved(configFilePath string, numWorkers int64, metricsAddr string, envFile string, loadDotEnv bool, subprocessArgs []string) error {
	// Drop any leading "--" separators that may be present after CLI parsing.
	for len(subprocessArgs) > 0 && subprocessArgs[0] == "--" {
		subprocessArgs = subprocessArgs[1:]
	}

	// Load environment files if requested
	if loadDotEnv {
		if err := loadEnvFile(".env"); err != nil {
			// Only warn if .env file exists but can't be read
			if !os.IsNotExist(err) {
				log.Printf("green-orb warning: failed to load .env file: %v", err)
			}
		}
	}

	if envFile != "" {
		if err := loadEnvFile(envFile); err != nil {
			log.Fatalf("green-orb error: failed to load env file %s: %v", envFile, err)
		}
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
    var lastWaitErr error
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
                _ = process_signal(observedCmd, sig)
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
        lastWaitErr = observedCmd.Wait()
        if lastWaitErr != nil {
            log.Printf("green-orb warning: observed process exited with error: %v", lastWaitErr)
        }

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
    if lastWaitErr != nil {
        var exitError *exec.ExitError
        if errors.As(lastWaitErr, &exitError) {
            os.Exit(exitError.ExitCode())
        }
        log.Fatalf("green-orb error: Error waiting for Command: %v", lastWaitErr)
    }

	return nil
}
