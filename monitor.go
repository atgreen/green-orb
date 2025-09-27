package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

// Monitor handles output monitoring from the observed process
type Monitor struct {
	pid             int
	compiledSignals []CompiledSignal
	workerPool      *WorkerPool
	channelMap      map[string]Channel
}

// NewMonitor creates a new output monitor
func NewMonitor(pid int, signals []CompiledSignal, pool *WorkerPool, channels map[string]Channel) *Monitor {
	return &Monitor{
		pid:             pid,
		compiledSignals: signals,
		workerPool:      pool,
		channelMap:      channels,
	}
}

// MonitorOutput processes output from stdout or stderr
func (m *Monitor) MonitorOutput(scanner *bufio.Scanner, isStderr bool) {
	// Increase buffer to accommodate long log lines (up to 10MB)
	const maxLogLine = 10 * 1024 * 1024
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, maxLogLine)

	stream := "stdout"
	if isStderr {
		stream = "stderr"
	}

	for scanner.Scan() {
		logLine := scanner.Text()

		// Update metrics
		if metricsEnable {
			orbEventsTotal.WithLabelValues(stream).Inc()
		}

		// Check if we should suppress this line
		shouldSuppress := false

		// Process each signal
		for _, signal := range m.compiledSignals {
			matches := signal.Regex.FindStringSubmatch(logLine)
			if matches != nil {
				channel, ok := m.channelMap[signal.Channel]
				if !ok {
					log.Printf("green-orb warning: signal references unknown channel %q", signal.Channel)
					continue
				}

				// Update metrics
				if metricsEnable {
					orbSignalsMatchedTotal.WithLabelValues(signal.Regex.String(), signal.Channel).Inc()
				}

				// Check if this is a suppress channel
				if channel.Type == "suppress" {
					shouldSuppress = true
				}

				// Queue the action
				req := ActionRequest{
					channel:   signal.Channel,
					timestamp: time.Now().Format(time.RFC3339),
					PID:       m.pid,
					logline:   logLine,
					matches:   matches,
				}

				if !m.workerPool.Enqueue(req) {
					log.Printf("green-orb warning: dropped action for channel %s (queue full or rate limited)", signal.Channel)
				}
			}
		}

		// Output the line if not suppressed
		if !shouldSuppress {
			if isStderr {
				fmt.Fprintln(getStderr(), logLine)
			} else {
				fmt.Fprintln(getStdout(), logLine)
			}
		}
	}

	if err := scanner.Err(); err != nil && err != io.EOF {
		log.Printf("green-orb error: Error reading from %s: %v", stream, err)
	}
}

// getStdout returns the stdout writer (abstracted for testing)
func getStdout() io.Writer {
	return os.Stdout
}

// getStderr returns the stderr writer (abstracted for testing)
func getStderr() io.Writer {
	return os.Stderr
}