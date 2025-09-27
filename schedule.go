package main

import (
    "fmt"
    "log"
    "time"

    cronlib "github.com/robfig/cron/v3"
)

// Schedule is an internal representation of a time-based signal
type Schedule struct {
    Name    string
    Channel string
    Every   string
    Cron    string
}

// BuildSchedulesFromSignals extracts schedules from signals
func BuildSchedulesFromSignals(signals []Signal) []Schedule {
    var out []Schedule
    for _, s := range signals {
        if s.Schedule == nil {
            continue
        }
        every := s.Schedule.Every
        cron := s.Schedule.Cron
        if every == "" && cron == "" {
            continue
        }
        out = append(out, Schedule{
            Name:    s.Name,
            Channel: s.Channel,
            Every:   every,
            Cron:    cron,
        })
    }
    return out
}

// ScheduleRunner manages time-based triggers
type ScheduleRunner struct {
    schedules []Schedule
    getPID    func() int
    channels  map[string]Channel
    workerPool *WorkerPool
    stopChan  chan struct{}
    tickers   []*time.Ticker
    cron      *cronlib.Cron
}

// NewScheduleRunner creates a new schedule runner
func NewScheduleRunner(schedules []Schedule, getPID func() int, channels map[string]Channel, pool *WorkerPool) *ScheduleRunner {
    return &ScheduleRunner{
        schedules: schedules,
        getPID:    getPID,
        channels:  channels,
        workerPool: pool,
        stopChan:  make(chan struct{}),
    }
}

// Start begins all configured schedules
func (sr *ScheduleRunner) Start() {
    // Initialize cron with 5-field standard plus optional seconds support
    sr.cron = cronlib.New(cronlib.WithParser(cronlib.NewParser(
        cronlib.SecondOptional | cronlib.Minute | cronlib.Hour | cronlib.Dom | cronlib.Month | cronlib.Dow | cronlib.Descriptor,
    )))

    for _, s := range sr.schedules {
        if s.Every != "" {
            interval, err := time.ParseDuration(s.Every)
            if err != nil {
                log.Printf("green-orb error: invalid schedule interval for %s: %v", s.Name, err)
                continue
            }
            ticker := time.NewTicker(interval)
            sr.tickers = append(sr.tickers, ticker)
            go func(sched Schedule, tk *time.Ticker) {
                for {
                    select {
                    case <-tk.C:
                        if signalManager.IsEnabled(sched.Name) {
                            sr.trigger(sched, "every")
                        }
                    case <-sr.stopChan:
                        tk.Stop()
                        return
                    }
                }
            }(s, ticker)
            continue
        }
        if s.Cron != "" {
            // Add cron job
            sched := s
            _, err := sr.cron.AddFunc(s.Cron, func() {
                if signalManager.IsEnabled(sched.Name) {
                    sr.trigger(sched, "cron")
                }
            })
            if err != nil {
                log.Printf("green-orb error: invalid cron spec for %s: %v", s.Name, err)
            }
        }
    }
    sr.cron.Start()
}

// Stop stops all schedules
func (sr *ScheduleRunner) Stop() {
    close(sr.stopChan)
    if sr.cron != nil {
        sr.cron.Stop()
    }
}

func (sr *ScheduleRunner) trigger(sched Schedule, kind string) {
    ch, ok := sr.channels[sched.Channel]
    if !ok {
        log.Printf("green-orb warning: schedule %s references unknown channel %s", sched.Name, sched.Channel)
        return
    }

    // Prepare a simple log line identifying the schedule
    msg := fmt.Sprintf("schedule '%s' tick", sched.Name)

    req := ActionRequest{
        channel:   ch.Name,
        timestamp: time.Now().Format(time.RFC3339),
        PID:       sr.getPID(),
        logline:   msg,
        matches:   []string{msg},
    }

    if !sr.workerPool.Enqueue(req) {
        log.Printf("green-orb warning: dropped schedule action for %s (queue full or rate limited)", sched.Name)
    }

    if metricsEnable {
        orbSchedulesFiredTotal.WithLabelValues(sched.Name, ch.Name, kind).Inc()
    }
}
