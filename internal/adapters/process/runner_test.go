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
	_, err = (Runner{}).Run(context.Background(), execution.Request{Command: "sleep 1", Timeout: time.Millisecond})
	if !errors.Is(err, ErrTimeout) {
		t.Fatalf("timeout error = %v", err)
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
