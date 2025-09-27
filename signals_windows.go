//go:build windows
// +build windows

// SPDX-License-Identifier: MIT
//
// signals_windows.go - dummy placeholder
//
// Copyright (C) 2023-2025 Anthony Green - green@moxielogic.com
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the “Software”), to deal in the Software without
// restriction, including without limitation the rights to use, copy,
// modify, merge, publish, distribute, sublicense, and/or sell copies
// of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
// BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
// ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

// process_signal handles signal forwarding on Windows
func process_signal(observed_cmd *exec.Cmd, sig os.Signal) error {
	if observed_cmd == nil || observed_cmd.Process == nil {
		return nil
	}

	// On Windows, we can only handle certain signals
	switch sig {
	case os.Interrupt:
		// CTRL+C event - send interrupt to the process
		return sendCtrlC(observed_cmd.Process.Pid)
	case os.Kill:
		// Kill signal - forcefully terminate the process
		return observed_cmd.Process.Kill()
	case syscall.SIGTERM:
		// Termination signal - gracefully terminate if possible
		return observed_cmd.Process.Signal(os.Kill)
	default:
		// Other signals are not supported on Windows
		log.Printf("green-orb warning: signal %v not supported on Windows", sig)
		return nil
	}
}

// sendCtrlC sends a CTRL+C event to a process on Windows
func sendCtrlC(pid int) error {
	// On Windows, we can use GenerateConsoleCtrlEvent to send CTRL+C
	// However, this requires the process to be in the same console group
	// For simplicity, we'll just kill the process
	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return proc.Signal(os.Interrupt)
}
