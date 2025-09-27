package main

import (
    "sync"
    "time"
)

// SignalManager manages runtime enable/disable state of signals
type SignalManager struct {
    mu     sync.Mutex
    state  map[string]bool
    timers map[string]*time.Timer
}

func NewSignalManager() *SignalManager {
    return &SignalManager{
        state:  make(map[string]bool),
        timers: make(map[string]*time.Timer),
    }
}

// Init sets initial states; missing names default to true
func (sm *SignalManager) Init(signals []Signal) {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    // Stop and clear any TTL timers from a previous run/restart
    for name, t := range sm.timers {
        if t != nil {
            t.Stop()
        }
        delete(sm.timers, name)
    }
    // Reset all known states; by default signals are enabled unless specified
    sm.state = make(map[string]bool)
    for _, s := range signals {
        // Only respect enabled for named signals; unnamed behave as default true
        if s.Name != "" && s.Enabled != nil {
            sm.state[s.Name] = *s.Enabled
        } else if s.Name != "" {
            sm.state[s.Name] = true
        }
    }
}

func (sm *SignalManager) IsEnabled(name string) bool {
    if name == "" {
        return true
    }
    sm.mu.Lock()
    defer sm.mu.Unlock()
    v, ok := sm.state[name]
    if !ok {
        return true
    }
    return v
}

// Enable sets enabled=true, optionally for a TTL after which it auto-disables
func (sm *SignalManager) Enable(name string, ttl *time.Duration) {
    if name == "" {
        return
    }
    sm.mu.Lock()
    defer sm.mu.Unlock()
    sm.state[name] = true
    if t, ok := sm.timers[name]; ok && t != nil {
        t.Stop()
        delete(sm.timers, name)
    }
    if ttl != nil && *ttl > 0 {
        d := *ttl
        sm.timers[name] = time.AfterFunc(d, func() {
            sm.Disable(name)
        })
    }
}

// Disable sets enabled=false and cancels any TTL timer
func (sm *SignalManager) Disable(name string) {
    if name == "" {
        return
    }
    sm.mu.Lock()
    defer sm.mu.Unlock()
    sm.state[name] = false
    if t, ok := sm.timers[name]; ok && t != nil {
        t.Stop()
        delete(sm.timers, name)
    }
}
