package process

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/nilstate/scafld-go/internal/core/execution"
)

var ErrTimeout = errors.New("process timeout")

type Runner struct {
	DiagnosticsDir string
}

func (r Runner) Run(ctx context.Context, req execution.Request) (execution.Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if req.Command == "" {
		return execution.Result{}, fmt.Errorf("command is required")
	}
	runCtx := ctx
	cancel := func() {}
	if req.Timeout > 0 {
		runCtx, cancel = context.WithTimeout(ctx, req.Timeout)
	}
	defer cancel()
	cmd := exec.CommandContext(runCtx, "sh", "-c", req.Command)
	if req.CWD != "" {
		cmd.Dir = req.CWD
	}
	cmd.Env = append(os.Environ(), req.Env...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	result := execution.Result{Output: out.String()}
	if cmd.ProcessState != nil {
		result.ExitCode = cmd.ProcessState.ExitCode()
	} else if err == nil {
		result.ExitCode = 0
	} else {
		result.ExitCode = 1
	}
	if r.DiagnosticsDir != "" {
		if path, writeErr := r.writeDiagnostic(req.Command, result.Output); writeErr == nil {
			result.DiagnosticPath = path
		}
	}
	if errors.Is(runCtx.Err(), context.DeadlineExceeded) {
		result.TimedOut = true
		return result, ErrTimeout
	}
	if errors.Is(runCtx.Err(), context.Canceled) {
		return result, runCtx.Err()
	}
	return result, nil
}

func (r Runner) writeDiagnostic(command string, output string) (string, error) {
	if err := os.MkdirAll(r.DiagnosticsDir, 0o755); err != nil {
		return "", err
	}
	name := fmt.Sprintf("command-%d.txt", time.Now().UnixNano())
	path := filepath.Join(r.DiagnosticsDir, name)
	data := []byte("command: " + command + "\n\n" + output)
	return path, os.WriteFile(path, data, 0o644)
}
