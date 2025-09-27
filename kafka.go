package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/plain"
)

// KafkaManager manages Kafka client connections
type KafkaManager struct {
	clients map[string]*kgo.Client
}

// NewKafkaManager creates a new Kafka manager
func NewKafkaManager() *KafkaManager {
	return &KafkaManager{
		clients: make(map[string]*kgo.Client),
	}
}

// Connect establishes connections to all configured Kafka brokers
func (km *KafkaManager) Connect(channels []Channel) error {
	for _, channel := range channels {
		if channel.Type != "kafka" {
			continue
		}

		client, err := km.createClient(channel)
		if err != nil {
			return fmt.Errorf("failed to create kafka client for channel %s: %w", channel.Name, err)
		}

		km.clients[channel.Name] = client
	}

	return nil
}

// GetClient returns a Kafka client for the specified channel
func (km *KafkaManager) GetClient(channelName string) (*kgo.Client, bool) {
	client, ok := km.clients[channelName]
	return client, ok
}

// Close closes all Kafka connections
func (km *KafkaManager) Close() {
	for name, client := range km.clients {
		if client != nil {
			client.Close()
		}
		delete(km.clients, name)
	}
}

// createClient creates a single Kafka client
func (km *KafkaManager) createClient(channel Channel) (*kgo.Client, error) {
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

	// Configure TLS if enabled
	if channel.TLSEnable {
		tlsCfg, err := km.buildTLSConfig(channel)
		if err != nil {
			return nil, err
		}
		opts = append(opts, kgo.DialTLSConfig(tlsCfg))
	}

	// Configure SASL authentication if enabled
	if channel.SASLMechanism == "plain" || (channel.SASLMechanism == "" && channel.SASLUsername != "") {
		mech := plain.Auth{
			User: channel.SASLUsername,
			Pass: channel.SASLPassword,
		}.AsMechanism()
		opts = append(opts, kgo.SASL(mech))
	}

	return kgo.NewClient(opts...)
}

// buildTLSConfig builds a TLS configuration for a channel
func (km *KafkaManager) buildTLSConfig(channel Channel) (*tls.Config, error) {
	tlsCfg := &tls.Config{
		InsecureSkipVerify: channel.TLSInsecureSkipVerify, //nolint:gosec // configurable for environments with self-signed certs
	}

	// Load CA certificate if provided
	if channel.TLSCAFile != "" {
		caPem, err := os.ReadFile(channel.TLSCAFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read TLS CA file: %w", err)
		}
		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(caPem) {
			return nil, fmt.Errorf("failed to parse TLS CA file")
		}
		tlsCfg.RootCAs = pool
	}

	// Load client certificate if provided
	if channel.TLSCertFile != "" && channel.TLSKeyFile != "" {
		certPem, err := os.ReadFile(channel.TLSCertFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read TLS cert file: %w", err)
		}
		keyPem, err := os.ReadFile(channel.TLSKeyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read TLS key file: %w", err)
		}
		cert, err := tls.X509KeyPair(certPem, keyPem)
		if err != nil {
			return nil, fmt.Errorf("failed to load TLS key pair: %w", err)
		}
		tlsCfg.Certificates = []tls.Certificate{cert}
	}

	return tlsCfg, nil
}