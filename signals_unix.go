//go:build linux || darwin
// +build linux darwin

// SPDX-License-Identifier: MIT
//
// signals_unix.go - pass signals through to observed command
//
// Copyright (C) 2023, 2024 Anthony Green - green@moxielogic.com
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

func process_signal(observed_cmd *exec.Cmd, sig os.Signal) error {
	if sig == syscall.SIGCHLD {
		return nil
	}
	if err := observed_cmd.Process.Signal(sig); err != nil {
		log.Println("green-orb error: Failed to send signal to subprocess:", err)
	}

	return nil
}
