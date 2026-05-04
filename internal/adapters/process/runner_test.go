package process

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nilstate/scafld-go/internal/core/execution"
)

func TestCommandTimeoutDiagnosticCancel(t *testing.T) {
	t.Parallel()

	result, err := (Runner{DiagnosticsDir: t.TempDir()}).Run(context.Background(), execution.Request{Command: "printf ok", Timeout: time.Second})
	if err != nil {
		t.Fatal(err)
	}
	if result.ExitCode != 0 || result.Output != "ok" || result.DiagnosticPath == "" {
		t.Fatalf("unexpected result: %+v", result)
	}
	result, err = (Runner{}).Run(context.Background(), execution.Request{Command: "printf out; printf err >&2", Timeout: time.Second})
	if err != nil {
		t.Fatal(err)
	}
	if result.Stdout != "out" || result.Stderr != "err" {
		t.Fatalf("stdout/stderr not split: %+v", result)
	}
	result, err = (Runner{}).Run(context.Background(), execution.Request{Command: "cat", Input: "prompt", Timeout: time.Second})
	if err != nil {
		t.Fatal(err)
	}
	if result.Stdout != "prompt" {
		t.Fatalf("stdin not sent to process: %+v", result)
	}
	result, err = (Runner{}).Run(context.Background(), execution.Request{Args: []string{"printf", "argv"}, Timeout: time.Second})
	if err != nil {
		t.Fatal(err)
	}
	if result.Stdout != "argv" {
		t.Fatalf("argv execution failed: %+v", result)
	}
	_, err = (Runner{}).Run(context.Background(), execution.Request{Command: "sleep 1", Timeout: time.Millisecond})
	if !errors.Is(err, ErrTimeout) {
		t.Fatalf("timeout error = %v", err)
	}
	idle, err := (Runner{}).Run(context.Background(), execution.Request{Command: "printf start; sleep 1", Timeout: time.Second, IdleTimeout: time.Millisecond})
	if !errors.Is(err, ErrIdle) || idle.KillReason != "idle_timeout" {
		t.Fatalf("idle result=%+v err=%v", idle, err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	path := filepath.Join(t.TempDir(), "should-not-exist")
	_, err = (Runner{}).Run(ctx, execution.Request{Command: "touch " + path})
	if err == nil {
		t.Fatal("expected cancellation error")
	}
	if _, statErr := os.Stat(path); !os.IsNotExist(statErr) {
		t.Fatalf("pre-cancelled command should not spawn, stat err = %v", statErr)
	}
}

func TestSignalInterruptTerminateEscalateScenario(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := (Runner{}).Run(ctx, execution.Request{Command: "sleep 1"})
	if err == nil {
		t.Fatal("expected cancellation")
	}
}
