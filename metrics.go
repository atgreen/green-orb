// SPDX-License-Identifier: MIT

package main

import (
    "log"
    "net/http"
    "time"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    orbEventsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "orb_events_total",
            Help: "Total log lines processed by stream",
        },
        []string{"stream"},
    )

    orbSignalsMatchedTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "orb_signals_matched_total",
            Help: "Total signal regex matches",
        },
        []string{"signal", "channel"},
    )

    orbActionsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "orb_actions_total",
            Help: "Total actions executed per channel and type",
        },
        []string{"channel", "type", "outcome"},
    )

    orbActionLatencySeconds = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "orb_action_latency_seconds",
            Help:    "Action execution latency in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"channel", "type"},
    )

    orbDroppedEventsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "orb_dropped_events_total",
            Help: "Total dropped events by reason",
        },
        []string{"reason"},
    )

    orbRestartsTotal = promauto.NewCounter(
        prometheus.CounterOpts{
            Name: "orb_restarts_total",
            Help: "Total restart actions triggered",
        },
    )

    orbQueueDepth = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "orb_queue_depth",
            Help: "Current notification queue depth",
        },
    )

    orbObservedPID = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "orb_observed_pid",
            Help: "PID of the observed process (0 when none)",
        },
    )

    orbChecksTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "orb_checks_total",
            Help: "Total checks executed by type and outcome",
        },
        []string{"type", "outcome"},
    )
)

func StartMetricsServer(addr string) {
    mux := http.NewServeMux()
    mux.Handle("/metrics", promhttp.Handler())
    srv := &http.Server{
        Addr:              addr,
        Handler:           mux,
        ReadHeaderTimeout: 5 * time.Second,
    }
    go func() {
        if err := srv.ListenAndServe(); err != nil {
            if err != http.ErrServerClosed {
                log.Println("green-orb warning: metrics server error:", err)
            }
        }
    }()
}

func StartQueueDepthGauge(ch chan Notification, stop <-chan struct{}) {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()
    for {
        select {
        case <-ticker.C:
            orbQueueDepth.Set(float64(len(ch)))
        case <-stop:
            return
        }
    }
}
