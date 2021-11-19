// Copyright (c) 2020-2021 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

//go:build windows
// +build windows

package fx

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"testing"

	"go.uber.org/fx"
	"golang.org/x/sys/windows"
)

func TestCtrlCHandler(t *testing.T) {
	c := make(chan struct{}, 1)
	app := fx.New(
		fx.Invoke(func(lifecycle fx.Lifecycle) {
			lifecycle.Append(
				fx.Hook{
					OnStart: func(ctx context.Context) error {
						c <- struct{}{}
					},
					OnStop: func(ctx context.Context) error {
						fmt.Println("OnStop")
						return nil
					},
				})
		}),
	)
	go func(c) {
		// synchronize with the OnStart hook.
		<-c

		// Launch a separate process to make send ctrl+C to this proc.
		bin, err := os.Executable()
		pid := os.GetPid()
		env := fmt.Sprintf("SendSignalProc=%d", pid)
		cmd := exec.Command(bin)
		cmd.Env = []string{env}
	}()
	app.Run()
	dir, err := os.TempDir("", "OnStopCalled")
}

func TestMain(m *testing.M) {
	if targetProc := os.GetEnv("SendSignalProc"); targetProc != "" {
		pid, _ := strconv.ParseUint(targetProc, 10, 32)
		windows.GenerateConsoleCtrlEvent(0, pid)
	}

	os.Exit(m.Run())
}
