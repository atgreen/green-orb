package main

import (
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
	"golang.org/x/time/rate"
)

// Global variables shared across modules
var (
	// Process control
	observedCmd   *exec.Cmd
	shouldRestart bool = true
	restartMutex  sync.Mutex

	// Restart tracking for flapping detection
	restartTimesMu sync.Mutex
	restartTimes   []time.Time

	// Kafka clients map
	kafkaClients map[string]*kgo.Client = make(map[string]*kgo.Client)

	// Channel rate limiters
	channelLimiters map[string]*rate.Limiter = make(map[string]*rate.Limiter)

	// Metrics enable flag
	metricsEnable bool

	// Channels map for global access
	channels map[string]Channel
)

// TemplateData contains data available to templates
type TemplateData struct {
	PID       int
	Logline   string
	Matches   []string
	Timestamp string
	Env       map[string]string
}

// Notification represents a legacy notification structure for compatibility
type Notification struct {
	PID     int
	Channel Channel
	Match   []string
	Message string
}

// envToMap converts environment variables to a map
func envToMap() map[string]string {
	envMap := make(map[string]string)
	for _, v := range os.Environ() {
		parts := strings.SplitN(v, "=", 2)
		if len(parts) == 2 {
			envMap[parts[0]] = parts[1]
		} else {
			envMap[parts[0]] = ""
		}
	}
	return envMap
}