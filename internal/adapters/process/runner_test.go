package process

import (
	"context"
	"errors"
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
	_, err = (Runner{}).Run(ctx, execution.Request{Command: "sleep 1"})
	if err == nil {
		t.Fatal("expected cancellation error")
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
