package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"text/template"
	"time"

	"github.com/nicholas-fedor/shoutrrr"
	"github.com/twmb/franz-go/pkg/kgo"
	"golang.org/x/time/rate"
)

// ActionRequest represents a queued action to be executed
type ActionRequest struct {
	channel   string
	timestamp string
	PID       int
	logline   string
	matches   []string
}


// WorkerPool manages action processing workers
type WorkerPool struct {
	queue          chan ActionRequest
	numWorkers     int
	wg             sync.WaitGroup
	channels       map[string]Channel
	kafkaClients   map[string]*kgo.Client
	channelLimiters map[string]*rate.Limiter
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(numWorkers int, queueSize int, channels map[string]Channel, kafkaClients map[string]*kgo.Client) *WorkerPool {
	pool := &WorkerPool{
		queue:          make(chan ActionRequest, queueSize),
		numWorkers:     numWorkers,
		channels:       channels,
		kafkaClients:   kafkaClients,
		channelLimiters: make(map[string]*rate.Limiter),
	}

	// Initialize rate limiters
	for name, ch := range channels {
		if ch.RatePerSec > 0 {
			burst := ch.Burst
			if burst <= 0 {
				burst = 1
			}
			pool.channelLimiters[name] = rate.NewLimiter(rate.Limit(ch.RatePerSec), burst)
		}
	}

	return pool
}

// Start begins the worker goroutines
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.numWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// Stop closes the queue and waits for workers to finish
func (wp *WorkerPool) Stop() {
	close(wp.queue)
	wp.wg.Wait()
}

// Enqueue attempts to add an action request to the queue
func (wp *WorkerPool) Enqueue(req ActionRequest) bool {
	// Check rate limiting
	if limiter, ok := wp.channelLimiters[req.channel]; ok {
		if !limiter.Allow() {
			if metricsEnable {
				orbDroppedEventsTotal.WithLabelValues("rate_limited").Inc()
			}
			return false
		}
	}

	// Non-blocking enqueue
	select {
	case wp.queue <- req:
		return true
	default:
		if metricsEnable {
			orbDroppedEventsTotal.WithLabelValues("queue_full").Inc()
		}
		return false
	}
}

// worker processes action requests from the queue
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	for req := range wp.queue {
		wp.processAction(req)
	}
}

// processAction executes a single action
func (wp *WorkerPool) processAction(req ActionRequest) {
	channel, ok := wp.channels[req.channel]
	if !ok {
		log.Printf("green-orb warning: unknown channel %s", req.channel)
		return
	}

	// Prepare template data
	env := envToMap()
	td := TemplateData{
		PID:       req.PID,
		Logline:   req.logline,
		Env:       env,
		Matches:   req.matches,
		Timestamp: req.timestamp,
	}

	start := time.Now()
	outcome := "success"

	switch channel.Type {
	case "notify":
		if err := wp.executeNotify(channel, td); err != nil {
			log.Printf("green-orb warning: notification failed: %v", err)
			outcome = "error"
		}
	case "exec":
		if err := wp.executeExec(channel, req.PID, req.matches); err != nil {
			log.Printf("green-orb warning: exec channel failed: %v", err)
			outcome = "error"
		}
	case "kafka":
		if err := wp.executeKafka(channel, req.logline); err != nil {
			log.Printf("green-orb warning: kafka send failed: %v", err)
			outcome = "error"
		}
	case "restart":
		wp.executeRestart()
	case "kill":
		wp.executeKill()
	case "suppress":
		// No action needed for suppress
	}

	// Update metrics
	if metricsEnable {
		orbActionsTotal.WithLabelValues(channel.Name, channel.Type, outcome).Inc()
		orbActionLatencySeconds.WithLabelValues(channel.Name, channel.Type).Observe(time.Since(start).Seconds())
	}
}

// executeNotify sends a notification
func (wp *WorkerPool) executeNotify(channel Channel, td TemplateData) error {
	// Process URL template
	urlTmpl, err := template.New("url").Parse(channel.URL)
	if err != nil {
		return fmt.Errorf("can't parse URL template: %w", err)
	}

	var urlBuf bytes.Buffer
	if err := urlTmpl.Execute(&urlBuf, td); err != nil {
		return fmt.Errorf("can't execute URL template: %w", err)
	}

	// Process message template
	var message string
	if channel.Template != "" {
		msgTmpl, err := template.New("msg").Parse(channel.Template)
		if err != nil {
			return fmt.Errorf("can't parse message template: %w", err)
		}
		var msgBuf bytes.Buffer
		if err := msgTmpl.Execute(&msgBuf, td); err != nil {
			return fmt.Errorf("can't execute message template: %w", err)
		}
		message = msgBuf.String()
	} else {
		message = td.Logline
	}

	return shoutrrr.Send(urlBuf.String(), message)
}

// executeExec runs a shell command
func (wp *WorkerPool) executeExec(channel Channel, pid int, matches []string) error {
	cmd := exec.Command("bash", "-c", channel.Shell)
	cmd.Env = makeExecEnvironment(pid, matches)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	return cmd.Run()
}

// executeKafka sends a message to Kafka
func (wp *WorkerPool) executeKafka(channel Channel, message string) error {
	client, ok := wp.kafkaClients[channel.Name]
	if !ok {
		return fmt.Errorf("kafka client not found for channel %s", channel.Name)
	}

	ctx := context.Background()
	record := &kgo.Record{
		Topic: channel.Topic,
		Value: []byte(message),
	}

	return client.ProduceSync(ctx, record).FirstErr()
}

// executeRestart signals the observed process to restart
func (wp *WorkerPool) executeRestart() {
	restartMutex.Lock()
	shouldRestart = true
	restartMutex.Unlock()

    if observedCmd != nil && observedCmd.Process != nil {
        _ = observedCmd.Process.Signal(syscall.SIGTERM)

		if metricsEnable {
			orbRestartsTotal.Inc()
		}

		restartTimesMu.Lock()
		restartTimes = append(restartTimes, time.Now())
		restartTimesMu.Unlock()
	}
}

// executeKill signals the observed process to terminate
func (wp *WorkerPool) executeKill() {
	restartMutex.Lock()
	shouldRestart = false
	restartMutex.Unlock()

    if observedCmd != nil && observedCmd.Process != nil {
        _ = observedCmd.Process.Signal(syscall.SIGTERM)
	}
}


// makeExecEnvironment prepares environment variables for exec channels
func makeExecEnvironment(pid int, matches []string) []string {
	env := os.Environ()
	env = append(env, fmt.Sprintf("ORB_PID=%d", pid))
	env = append(env, fmt.Sprintf("ORB_MATCH_COUNT=%d", len(matches)))
	for i, match := range matches {
		env = append(env, fmt.Sprintf("ORB_MATCH_%d=%s", i, match))
	}
	return env
}
