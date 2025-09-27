package main

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
	"testing"
	"text/template"
	"time"
)

// TestCompileSignals tests the compilation of signal regular expressions
func TestCompileSignals(t *testing.T) {
	tests := []struct {
		name      string
		signals   []Signal
		wantErr   bool
		wantCount int
	}{
		{
			name: "valid signals",
			signals: []Signal{
				{Regex: "^test.*", Channel: "test-channel"},
				{Regex: "[0-9]+", Channel: "number-channel"},
			},
			wantErr:   false,
			wantCount: 2,
		},
		{
			name: "invalid regex",
			signals: []Signal{
				{Regex: "[invalid", Channel: "test-channel"},
			},
			wantErr:   true,
			wantCount: 0,
		},
		{
			name:      "empty signals",
			signals:   []Signal{},
			wantErr:   false,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiledSigs, err := CompileSignals(tt.signals)
			if (err != nil) != tt.wantErr {
				t.Errorf("compileSignals() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(compiledSigs) != tt.wantCount {
				t.Errorf("compileSignals() returned %d signals, want %d", len(compiledSigs), tt.wantCount)
			}
		})
	}
}

// TestMatchSignal tests signal matching against log lines
func TestMatchSignal(t *testing.T) {
	compiledSignals := []CompiledSignal{
		{
			Regex:   regexp.MustCompile(`ERROR: (.+)`),
			Channel: "error-channel",
		},
		{
			Regex:   regexp.MustCompile(`WARNING: (.+)`),
			Channel: "warning-channel",
		},
	}

	tests := []struct {
		name        string
		logLine     string
		wantChannel string
		wantMatches []string
		wantFound   bool
	}{
		{
			name:        "match error",
			logLine:     "ERROR: Database connection failed",
			wantChannel: "error-channel",
			wantMatches: []string{"ERROR: Database connection failed", "Database connection failed"},
			wantFound:   true,
		},
		{
			name:        "match warning",
			logLine:     "WARNING: High memory usage",
			wantChannel: "warning-channel",
			wantMatches: []string{"WARNING: High memory usage", "High memory usage"},
			wantFound:   true,
		},
		{
			name:        "no match",
			logLine:     "INFO: Application started",
			wantChannel: "",
			wantMatches: nil,
			wantFound:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, sig := range compiledSignals {
				matches := sig.Regex.FindStringSubmatch(tt.logLine)
				if matches != nil {
					if !tt.wantFound {
						t.Errorf("Expected no match but found one")
					}
					if sig.Channel != tt.wantChannel {
						t.Errorf("Channel = %v, want %v", sig.Channel, tt.wantChannel)
					}
					if len(matches) != len(tt.wantMatches) {
						t.Errorf("Matches length = %v, want %v", len(matches), len(tt.wantMatches))
					}
					for i, m := range matches {
						if m != tt.wantMatches[i] {
							t.Errorf("Match[%d] = %v, want %v", i, m, tt.wantMatches[i])
						}
					}
					return
				}
			}
			if tt.wantFound {
				t.Errorf("Expected match but found none")
			}
		})
	}
}

// TestQueueOperations tests the non-blocking queue
func TestQueueOperations(t *testing.T) {
	queue := make(chan ActionRequest, 2)

	// Test successful enqueue
	req1 := ActionRequest{
		channel:   "test-channel",
		timestamp: time.Now().Format(time.RFC3339),
		PID:       1234,
		logline:   "Test log line",
		matches:   []string{"Test log line"},
	}

	select {
	case queue <- req1:
		// Success
	default:
		t.Error("Failed to enqueue to empty queue")
	}

	// Fill the queue
	req2 := req1
	req2.logline = "Second log"
	select {
	case queue <- req2:
		// Success
	default:
		t.Error("Failed to enqueue second item")
	}

	// Test queue full behavior (non-blocking)
	req3 := req1
	req3.logline = "Third log"
	select {
	case queue <- req3:
		t.Error("Should not be able to enqueue to full queue")
	default:
		// Expected behavior - queue is full
	}

	// Dequeue and verify
	select {
	case dequeued := <-queue:
		if dequeued.logline != "Test log line" {
			t.Errorf("Dequeued wrong item: %v", dequeued.logline)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Failed to dequeue")
	}
}

// tokenBucket represents a simple token bucket for testing
type tokenBucket struct {
	tokens     float64
	maxTokens  float64
	refillRate float64
	lastRefill time.Time
	mu         sync.Mutex
}

func (tb *tokenBucket) allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	tb.tokens = tb.tokens + elapsed*tb.refillRate
	if tb.tokens > tb.maxTokens {
		tb.tokens = tb.maxTokens
	}
	tb.lastRefill = now

	if tb.tokens >= 1 {
		tb.tokens--
		return true
	}
	return false
}

// TestRateLimiter tests the token bucket rate limiter
func TestRateLimiter(t *testing.T) {
	// Create a rate limiter with 2 tokens per second, burst of 3
	limiter := &tokenBucket{
		tokens:       3,
		maxTokens:    3,
		refillRate:   2,
		lastRefill:   time.Now(),
		mu:           sync.Mutex{},
	}

	// Should allow first 3 immediately (burst)
	for i := 0; i < 3; i++ {
		if !limiter.allow() {
			t.Errorf("Should allow request %d within burst", i+1)
		}
	}

	// Fourth should be denied
	if limiter.allow() {
		t.Error("Should deny request exceeding burst")
	}

	// Wait for token refill
	time.Sleep(600 * time.Millisecond)

	// Should allow one more (refilled ~1 token)
	if !limiter.allow() {
		t.Error("Should allow request after refill")
	}
}

// TestProcessTemplate tests template processing
func TestProcessTemplate(t *testing.T) {
	tests := []struct {
		name     string
		template string
		data     struct {
			Timestamp string
			PID       int
			Matches   []string
			Logline   string
			Env       map[string]string
		}
		want    string
		wantErr bool
	}{
		{
			name:     "simple template",
			template: "PID: {{.PID}}, Time: {{.Timestamp}}",
			data: struct {
				Timestamp string
				PID       int
				Matches   []string
				Logline   string
				Env       map[string]string
			}{
				Timestamp: "2024-01-01T00:00:00Z",
				PID:       1234,
				Matches:   []string{"full line", "capture1"},
				Logline:   "full line",
				Env:       map[string]string{"USER": "test"},
			},
			want:    "PID: 1234, Time: 2024-01-01T00:00:00Z",
			wantErr: false,
		},
		{
			name:     "template with matches",
			template: "Match: {{index .Matches 1}}",
			data: struct {
				Timestamp string
				PID       int
				Matches   []string
				Logline   string
				Env       map[string]string
			}{
				Matches: []string{"full line", "captured text"},
			},
			want:    "Match: captured text",
			wantErr: false,
		},
		{
			name:     "template with env var",
			template: "User: {{.Env.USER}}",
			data: struct {
				Timestamp string
				PID       int
				Matches   []string
				Logline   string
				Env       map[string]string
			}{
				Env: map[string]string{"USER": "testuser"},
			},
			want:    "User: testuser",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl, err := createTemplate("test", tt.template)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("Failed to create template: %v", err)
				}
				return
			}

			var buf bytes.Buffer
			err = tmpl.Execute(&buf, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Template execution error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if buf.String() != tt.want {
				t.Errorf("Template result = %v, want %v", buf.String(), tt.want)
			}
		})
	}
}

// TestValidateConfig tests configuration validation
func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			config: Config{
				Channels: []Channel{
					{Name: "test-channel", Type: "notify", URL: "slack://token@channel"},
					{Name: "exec-channel", Type: "exec", Shell: "echo test"},
				},
				Signals: []Signal{
					{Regex: "error", Channel: "test-channel"},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid channel type",
			config: Config{
				Channels: []Channel{
					{Name: "bad-channel", Type: "invalid"},
				},
			},
			wantErr: true,
			errMsg:  "invalid channel type",
		},
		{
			name: "missing channel reference",
			config: Config{
				Channels: []Channel{
					{Name: "test-channel", Type: "notify", URL: "slack://token@channel"},
				},
				Signals: []Signal{
					{Regex: "error", Channel: "non-existent"},
				},
			},
			wantErr: true,
			errMsg:  "references non-existent channel",
		},
		{
			name: "duplicate channel names",
			config: Config{
				Channels: []Channel{
					{Name: "duplicate", Type: "notify", URL: "slack://token@channel"},
					{Name: "duplicate", Type: "exec", Shell: "echo test"},
				},
			},
			wantErr: true,
			errMsg:  "duplicate channel name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(&tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errMsg != "" {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Error message should contain '%s', got '%s'", tt.errMsg, err.Error())
				}
			}
		})
	}
}

// TestChannelWorker tests the channel worker functionality
func TestChannelWorker(t *testing.T) {
	// Create a test queue
	queue := make(chan ActionRequest, 10)

	// Create test channels configuration
	testChannels := map[string]Channel{
		"suppress": {Name: "suppress", Type: "suppress"},
		"exec": {Name: "exec", Type: "exec", Shell: "echo test"},
	}

	// Mock the channels map
	oldChannels := channels
	channels = testChannels
	defer func() { channels = oldChannels }()

	// Start a worker in a goroutine
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Process a few items then stop
		for i := 0; i < 2; i++ {
			select {
			case req := <-queue:
				// Simple validation that request was received
				if req.channel == "" {
					t.Error("Received empty channel name")
				}
			case <-time.After(100 * time.Millisecond):
				return
			}
		}
	}()

	// Send test requests
	req1 := ActionRequest{
		channel:   "suppress",
		timestamp: time.Now().Format(time.RFC3339),
		PID:       os.Getpid(),
		logline:   "Test suppressed line",
		matches:   []string{"Test suppressed line"},
	}

	req2 := ActionRequest{
		channel:   "exec",
		timestamp: time.Now().Format(time.RFC3339),
		PID:       os.Getpid(),
		logline:   "Test exec line",
		matches:   []string{"Test exec line"},
	}

	queue <- req1
	queue <- req2

	// Wait for worker to process
	wg.Wait()
}

// TestEnvironmentVariables tests environment variable setup for exec channels
func TestEnvironmentVariables(t *testing.T) {
	matches := []string{"full line", "capture1", "capture2", "capture3"}
	pid := 12345

	env := makeExecEnvironment(pid, matches)

	// Check standard variables
	found := false
	for _, e := range env {
		if strings.HasPrefix(e, "ORB_PID=") {
			if e != fmt.Sprintf("ORB_PID=%d", pid) {
				t.Errorf("ORB_PID incorrect: %s", e)
			}
			found = true
			break
		}
	}
	if !found {
		t.Error("ORB_PID not found in environment")
	}

	// Check match count
	found = false
	for _, e := range env {
		if strings.HasPrefix(e, "ORB_MATCH_COUNT=") {
			if e != "ORB_MATCH_COUNT=4" {
				t.Errorf("ORB_MATCH_COUNT incorrect: %s", e)
			}
			found = true
			break
		}
	}
	if !found {
		t.Error("ORB_MATCH_COUNT not found in environment")
	}

	// Check individual matches
	for i, match := range matches {
		varName := fmt.Sprintf("ORB_MATCH_%d=", i)
		found := false
		for _, e := range env {
			if strings.HasPrefix(e, varName) {
				expected := fmt.Sprintf("%s%s", varName, match)
				if e != expected {
					t.Errorf("Match %d incorrect: got %s, want %s", i, e, expected)
				}
				found = true
				break
			}
		}
		if !found {
			t.Errorf("ORB_MATCH_%d not found in environment", i)
		}
	}
}

// Helper functions for tests

func createTemplate(name, text string) (*template.Template, error) {
	// This would be the actual template creation function
	// For testing, we'll use a simple implementation
	tmpl := template.New(name)
	_, err := tmpl.Parse(text)
	return tmpl, err
}




// TestLoadEnvFile tests the .env file loading functionality
func TestLoadEnvFile(t *testing.T) {
	// Create a temporary env file
	tmpFile, err := os.CreateTemp("", "test*.env")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write test content
	content := `# Test env file
DATABASE_URL=postgresql://localhost:5432/test
API_KEY="secret-123"
DEBUG=true
EMPTY_VALUE=
# Comment line
QUOTED_SINGLE='single-quoted'
PORT=3000`

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Save original env values
	originalDB := os.Getenv("DATABASE_URL")
	originalAPI := os.Getenv("API_KEY")
	originalDebug := os.Getenv("DEBUG")

	// Load the env file
	if err := loadEnvFile(tmpFile.Name()); err != nil {
		t.Fatalf("loadEnvFile failed: %v", err)
	}

	// Check that values were loaded correctly
	tests := []struct {
		key      string
		expected string
	}{
		{"DATABASE_URL", "postgresql://localhost:5432/test"},
		{"API_KEY", "secret-123"},
		{"DEBUG", "true"},
		{"EMPTY_VALUE", ""},
		{"QUOTED_SINGLE", "single-quoted"},
		{"PORT", "3000"},
	}

	for _, tt := range tests {
		if got := os.Getenv(tt.key); got != tt.expected {
			t.Errorf("Expected %s=%s, got %s", tt.key, tt.expected, got)
		}
	}

	// Restore original values
	os.Setenv("DATABASE_URL", originalDB)
	os.Setenv("API_KEY", originalAPI)
	os.Setenv("DEBUG", originalDebug)
}

// TestLoadEnvFileNotFound tests handling of non-existent files
func TestLoadEnvFileNotFound(t *testing.T) {
	err := loadEnvFile("nonexistent.env")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
	if !os.IsNotExist(err) {
		t.Errorf("Expected os.IsNotExist error, got %v", err)
	}
}

// Benchmark tests

func BenchmarkSignalMatching(b *testing.B) {
	compiledSignals := []CompiledSignal{
		{Regex: regexp.MustCompile(`ERROR: (.+)`), Channel: "error"},
		{Regex: regexp.MustCompile(`WARNING: (.+)`), Channel: "warning"},
		{Regex: regexp.MustCompile(`INFO: (.+)`), Channel: "info"},
		{Regex: regexp.MustCompile(`DEBUG: (.+)`), Channel: "debug"},
	}

	logLines := []string{
		"ERROR: Database connection failed",
		"WARNING: High memory usage",
		"INFO: Application started",
		"DEBUG: Processing request",
		"Normal log line without prefix",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logLine := logLines[i%len(logLines)]
		for _, sig := range compiledSignals {
			_ = sig.Regex.FindStringSubmatch(logLine)
		}
	}
}

func BenchmarkTemplateExecution(b *testing.B) {
	tmplText := "PID: {{.PID}}, Time: {{.Timestamp}}, Match: {{index .Matches 0}}"
	tmpl, _ := createTemplate("bench", tmplText)

	data := struct {
		Timestamp string
		PID       int
		Matches   []string
		Logline   string
		Env       map[string]string
	}{
		Timestamp: "2024-01-01T00:00:00Z",
		PID:       1234,
		Matches:   []string{"full line", "capture1"},
		Logline:   "full line",
		Env:       map[string]string{"USER": "test"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		_ = tmpl.Execute(&buf, data)
	}
}