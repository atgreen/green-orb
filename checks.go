package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"regexp"
	"sync"
	"time"
)

// CheckScheduler manages periodic health checks
type CheckScheduler struct {
	checks     []Check
	getPID     func() int
	channels   map[string]Channel
	workerPool *WorkerPool
	stopChan   chan struct{}
	wg         sync.WaitGroup
}

// NewCheckScheduler creates a new check scheduler
func NewCheckScheduler(checks []Check, getPID func() int, channels map[string]Channel, pool *WorkerPool) *CheckScheduler {
	return &CheckScheduler{
		checks:     checks,
		getPID:     getPID,
		channels:   channels,
		workerPool: pool,
		stopChan:   make(chan struct{}),
	}
}

// Start begins running scheduled checks
func (cs *CheckScheduler) Start() {
	for _, check := range cs.checks {
		cs.wg.Add(1)
		go cs.runCheck(check)
	}
}

// Stop halts all scheduled checks
func (cs *CheckScheduler) Stop() {
	close(cs.stopChan)
	cs.wg.Wait()
}

// runCheck runs a single check on its interval
func (cs *CheckScheduler) runCheck(check Check) {
	defer cs.wg.Done()

	// Parse interval duration
	interval, err := time.ParseDuration(check.Interval)
	if err != nil {
		log.Printf("green-orb error: invalid interval for check %s: %v", check.Name, err)
		return
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cs.executeCheck(check)
		case <-cs.stopChan:
			return
		}
	}
}

// executeCheck executes a single check
func (cs *CheckScheduler) executeCheck(check Check) {
	var err error

	switch check.Type {
	case "http":
		err = cs.executeHTTPCheck(check)
	case "tcp":
		err = cs.executeTCPCheck(check)
	case "flapping":
		err = cs.executeFlappingCheck(check)
	default:
		log.Printf("green-orb warning: unknown check type %s", check.Type)
		return
	}

	// Update metrics
	if metricsEnable {
		outcome := "success"
		if err != nil {
			outcome = "error"
		}
		orbChecksTotal.WithLabelValues(check.Type, outcome).Inc()
	}

	// Send notification on failure
	if err != nil {
		message := fmt.Sprintf("Check '%s' failed: %v", check.Name, err)
		cs.sendCheckNotification(check, message)
	}
}

// executeHTTPCheck performs an HTTP health check
func (cs *CheckScheduler) executeHTTPCheck(check Check) error {
	// Parse timeout
	timeout := 5 * time.Second
	if check.Timeout != "" {
		parsed, err := time.ParseDuration(check.Timeout)
		if err == nil {
			timeout = parsed
		}
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: timeout,
	}

	// Make request
	resp, err := client.Get(check.URL)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	expectedStatus := check.ExpectStatus
	if expectedStatus == 0 {
		expectedStatus = 200
	}
	if resp.StatusCode != expectedStatus {
		return fmt.Errorf("unexpected status code: got %d, want %d", resp.StatusCode, expectedStatus)
	}

	// Check body regex if specified
	if check.BodyRegex != "" {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		matched, err := regexp.MatchString(check.BodyRegex, string(body))
		if err != nil {
			return fmt.Errorf("invalid body regex: %w", err)
		}
		if !matched {
			return fmt.Errorf("body does not match expected pattern")
		}
	}

	return nil
}

// executeTCPCheck performs a TCP connectivity check
func (cs *CheckScheduler) executeTCPCheck(check Check) error {
	// Parse timeout
	timeout := 3 * time.Second
	if check.Timeout != "" {
		parsed, err := time.ParseDuration(check.Timeout)
		if err == nil {
			timeout = parsed
		}
	}

    // Attempt connection (use JoinHostPort for IPv4/IPv6 compatibility)
    address := net.JoinHostPort(check.Host, strconv.Itoa(check.Port))
    conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return fmt.Errorf("TCP connection failed: %w", err)
	}
	conn.Close()

	return nil
}

// executeFlappingCheck checks if the process is restarting too frequently
func (cs *CheckScheduler) executeFlappingCheck(check Check) error {
	// Parse window duration
	window := 5 * time.Minute
	if check.Window != "" {
		parsed, err := time.ParseDuration(check.Window)
		if err == nil {
			window = parsed
		}
	}

	// Count recent restarts
	restartTimesMu.Lock()
	defer restartTimesMu.Unlock()

	cutoff := time.Now().Add(-window)
	recentRestarts := 0

	// Count restarts within the window
	for _, restartTime := range restartTimes {
		if restartTime.After(cutoff) {
			recentRestarts++
		}
	}

	// Clean up old restart times
	var cleaned []time.Time
	for _, restartTime := range restartTimes {
		if restartTime.After(cutoff) {
			cleaned = append(cleaned, restartTime)
		}
	}
	restartTimes = cleaned

	// Check threshold
	if recentRestarts >= check.RestartThreshold {
		return fmt.Errorf("process restarted %d times in %v (threshold: %d)",
			recentRestarts, window, check.RestartThreshold)
	}

	return nil
}

// sendCheckNotification sends a notification for a failed check
func (cs *CheckScheduler) sendCheckNotification(check Check, message string) {
	channel, ok := cs.channels[check.Channel]
	if !ok {
		log.Printf("green-orb warning: check %s references unknown channel %s", check.Name, check.Channel)
		return
	}

	// Skip if channel is not a notification type
	if channel.Type != "notify" && channel.Type != "exec" && channel.Type != "kafka" {
		return
	}

	req := ActionRequest{
		channel:   check.Channel,
		timestamp: time.Now().Format(time.RFC3339),
		PID:       cs.getPID(),
		logline:   message,
		matches:   []string{message},
	}

	if !cs.workerPool.Enqueue(req) {
		log.Printf("green-orb warning: dropped check notification for %s (queue full or rate limited)", check.Name)
	}
}

// startChecksScheduler creates and starts a check scheduler (compatibility function)
func startChecksScheduler(checks []Check, getPID func() int, channels map[string]Channel, queue chan Notification) func() {
	if len(checks) == 0 {
		return func() {}
	}

	// Create a temporary worker pool adapter for compatibility
	// This will be removed when we fully refactor main.go
	tempPool := &WorkerPool{
		queue: make(chan ActionRequest, 100),
	}

	scheduler := NewCheckScheduler(checks, getPID, channels, tempPool)
	scheduler.Start()

	return scheduler.Stop
}
