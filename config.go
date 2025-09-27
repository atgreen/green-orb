package main

import (
    "fmt"
    "os"
    "regexp"
    "time"

    "gopkg.in/yaml.v3"
)

// Config represents the complete green-orb configuration
type Config struct {
    Channels []Channel `yaml:"channels"`
    Signals  []Signal  `yaml:"signals"`
    Checks   []Check   `yaml:"checks"`
}

// Channel represents a notification/action channel
type Channel struct {
	Name     string `yaml:"name"`
	Type     string `yaml:"type"`
	URL      string `yaml:"url"`
	Template string `yaml:"template"`
	Topic    string `yaml:"topic"`
	Broker   string `yaml:"broker"`
	Shell    string `yaml:"shell"`
	// Kafka auth/TLS (optional)
	SASLMechanism         string `yaml:"sasl_mechanism"`
	SASLUsername          string `yaml:"sasl_username"`
	SASLPassword          string `yaml:"sasl_password"`
	TLSEnable             bool   `yaml:"tls"`
	TLSInsecureSkipVerify bool   `yaml:"tls_insecure_skip_verify"`
	TLSCAFile             string `yaml:"tls_ca_file"`
	TLSCertFile           string `yaml:"tls_cert_file"`
	TLSKeyFile            string `yaml:"tls_key_file"`
	// Rate limiting (optional)
	RatePerSec float64 `yaml:"rate_per_sec"`
	Burst      int     `yaml:"burst"`
}

// Signal represents a log pattern to channel mapping
type Signal struct {
    Name     string        `yaml:"name,omitempty"`
    Regex    string        `yaml:"regex,omitempty"`
    Channel  string        `yaml:"channel"`
    Schedule *ScheduleSpec `yaml:"schedule,omitempty"`
}

// ScheduleSpec describes the time-based trigger for a signal
type ScheduleSpec struct {
    Every string `yaml:"every,omitempty"`
    Cron  string `yaml:"cron,omitempty"`
}

// (no top-level schedules; schedule is part of Signal)

// Check represents a periodic health check configuration
type Check struct {
	Name     string `yaml:"name"`
	Type     string `yaml:"type"` // http, tcp, flapping
	Channel  string `yaml:"channel"`
	Interval string `yaml:"interval"`
	Timeout  string `yaml:"timeout"`
	// HTTP-specific
	URL          string `yaml:"url"`
	ExpectStatus int    `yaml:"expect_status"`
	BodyRegex    string `yaml:"body_regex"`
	// TCP-specific
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	// Flapping-specific
	RestartThreshold int    `yaml:"restart_threshold"`
	Window           string `yaml:"window"`
}

// CompiledSignal represents a signal with compiled regex
type CompiledSignal struct {
	Regex   *regexp.Regexp
	Channel string
}

// LoadConfig loads and validates the configuration from a YAML file
func LoadConfig(filename string) (*Config, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(bytes, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// validateConfig validates the configuration
func validateConfig(config *Config) error {
	// Create a map of channel names for validation
	channelMap := make(map[string]bool)

	// Validate channels
	for _, ch := range config.Channels {
		if ch.Name == "" {
			return fmt.Errorf("channel missing required 'name' field")
		}
		if channelMap[ch.Name] {
			return fmt.Errorf("duplicate channel name: %s", ch.Name)
		}
		channelMap[ch.Name] = true

		// Validate channel type
		validTypes := []string{"notify", "kafka", "exec", "suppress", "restart", "kill"}
		validType := false
		for _, t := range validTypes {
			if ch.Type == t {
				validType = true
				break
			}
		}
		if !validType {
			return fmt.Errorf("invalid channel type '%s' for channel '%s'", ch.Type, ch.Name)
		}

		// Type-specific validation
		switch ch.Type {
		case "notify":
			if ch.URL == "" {
				return fmt.Errorf("notify channel '%s' missing required 'url' field", ch.Name)
			}
		case "kafka":
			if ch.Broker == "" {
				return fmt.Errorf("kafka channel '%s' missing required 'broker' field", ch.Name)
			}
			if ch.Topic == "" {
				return fmt.Errorf("kafka channel '%s' missing required 'topic' field", ch.Name)
			}
		case "exec":
			if ch.Shell == "" {
				return fmt.Errorf("exec channel '%s' missing required 'shell' field", ch.Name)
			}
		}
	}

    // Validate signals (regex or schedule)
    for i, sig := range config.Signals {
        if sig.Channel == "" {
            return fmt.Errorf("signal %d missing required 'channel' field", i)
        }
        if !channelMap[sig.Channel] {
            return fmt.Errorf("signal references non-existent channel: %s", sig.Channel)
        }

        hasRegex := sig.Regex != ""
        hasSchedule := sig.Schedule != nil && (sig.Schedule.Every != "" || sig.Schedule.Cron != "")
        if !hasRegex && !hasSchedule {
            return fmt.Errorf("signal %d must specify either 'regex' or 'schedule'", i)
        }

        if hasRegex {
            if _, err := regexp.Compile(sig.Regex); err != nil {
                return fmt.Errorf("invalid regex in signal %d: %w", i, err)
            }
        }
        if hasSchedule {
            // Require a name for schedule signals
            if sig.Name == "" {
                return fmt.Errorf("signal %d with 'schedule' must have a 'name'", i)
            }
            if sig.Schedule.Every != "" && sig.Schedule.Cron != "" {
                return fmt.Errorf("signal %d schedule must specify only one of 'every' or 'cron'", i)
            }
            if sig.Schedule.Every != "" {
                if _, err := time.ParseDuration(sig.Schedule.Every); err != nil {
                    return fmt.Errorf("signal %d has invalid schedule.every: %w", i, err)
                }
            }
            if sig.Schedule.Cron != "" {
                // Basic presence validation; detailed cron spec is validated at runtime by cron parser
            }
        }
    }

    // Validate checks
    for _, check := range config.Checks {
        if check.Name == "" {
            return fmt.Errorf("check missing required 'name' field")
        }
        if check.Channel == "" {
            return fmt.Errorf("check '%s' missing required 'channel' field", check.Name)
        }
        if !channelMap[check.Channel] {
            return fmt.Errorf("check '%s' references non-existent channel: %s", check.Name, check.Channel)
        }

        // Type-specific validation
        switch check.Type {
        case "http":
            if check.URL == "" {
                return fmt.Errorf("http check '%s' missing required 'url' field", check.Name)
            }
        case "tcp":
            if check.Host == "" {
                return fmt.Errorf("tcp check '%s' missing required 'host' field", check.Name)
            }
            if check.Port == 0 {
                return fmt.Errorf("tcp check '%s' missing required 'port' field", check.Name)
            }
        case "flapping":
            if check.RestartThreshold <= 0 {
                return fmt.Errorf("flapping check '%s' missing or invalid 'restart_threshold' field", check.Name)
            }
        default:
            return fmt.Errorf("invalid check type '%s' for check '%s'", check.Type, check.Name)
        }
    }

	return nil
}

// CompileSignals compiles all signal regex patterns
func CompileSignals(signals []Signal) ([]CompiledSignal, error) {
    var compiledSignals []CompiledSignal
    for _, signal := range signals {
        if signal.Regex == "" {
            // skip schedule-only signals here
            continue
        }
        re, err := regexp.Compile(signal.Regex)
        if err != nil {
            return nil, fmt.Errorf("failed to compile regex '%s': %w", signal.Regex, err)
        }
        compiledSignals = append(compiledSignals, CompiledSignal{
            Regex:   re,
            Channel: signal.Channel,
        })
    }
    return compiledSignals, nil
}

// CreateChannelMap creates a map from channel names to channels
func CreateChannelMap(channels []Channel) map[string]Channel {
	channelMap := make(map[string]Channel)
	for _, ch := range channels {
		channelMap[ch.Name] = ch
	}
	return channelMap
}
