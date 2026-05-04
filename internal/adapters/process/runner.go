package process

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/nilstate/scafld-go/internal/core/execution"
)

var (
	ErrTimeout = errors.New("process timeout")
	ErrIdle    = errors.New("process idle timeout")
)

type Runner struct {
	DiagnosticsDir string
}

func (r Runner) Run(ctx context.Context, req execution.Request) (execution.Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if req.Command == "" && len(req.Args) == 0 {
		return execution.Result{}, fmt.Errorf("command or args are required")
	}
	if req.Command != "" && len(req.Args) > 0 {
		return execution.Result{}, fmt.Errorf("command and args are mutually exclusive")
	}
	if err := ctx.Err(); err != nil {
		return execution.Result{ExitCode: -1, KillReason: "cancelled"}, err
	}
	displayCommand := req.Command
	var cmd *exec.Cmd
	if len(req.Args) > 0 {
		if strings.TrimSpace(req.Args[0]) == "" {
			return execution.Result{}, fmt.Errorf("args[0] is required")
		}
		cmd = exec.Command(req.Args[0], req.Args[1:]...)
		displayCommand = strings.Join(req.Args, " ")
	} else {
		cmd = exec.Command("sh", "-c", req.Command)
	}
	if req.CWD != "" {
		cmd.Dir = req.CWD
	}
	if req.Input != "" {
		cmd.Stdin = strings.NewReader(req.Input)
	}
	cmd.Env = append(os.Environ(), req.Env...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return execution.Result{}, fmt.Errorf("stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return execution.Result{}, fmt.Errorf("stderr pipe: %w", err)
	}
	if err := cmd.Start(); err != nil {
		return execution.Result{}, fmt.Errorf("start command: %w", err)
	}
	state := newCapture()
	state.touch()
	var pumps sync.WaitGroup
	pumps.Add(2)
	go pump(&pumps, stdout, state, "stdout")
	go pump(&pumps, stderr, state, "stderr")
	waitCh := make(chan error, 1)
	go func() {
		err := cmd.Wait()
		pumps.Wait()
		waitCh <- err
	}()
	result, waitErr := waitProcess(ctx, cmd, waitCh, state, req)
	result.Stdout, result.Stderr, result.Output = state.snapshot()
	if r.DiagnosticsDir != "" {
		if path, writeErr := r.writeDiagnostic(displayCommand, result); writeErr == nil {
			result.DiagnosticPath = path
		}
	}
	return result, waitErr
}

type capture struct {
	mu       sync.Mutex
	stdout   bytes.Buffer
	stderr   bytes.Buffer
	combined bytes.Buffer
	last     time.Time
}

func newCapture() *capture {
	return &capture{last: time.Now()}
}

func (c *capture) touch() {
	c.mu.Lock()
	c.last = time.Now()
	c.mu.Unlock()
}

func (c *capture) write(stream string, data []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	switch stream {
	case "stdout":
		c.stdout.Write(data)
	case "stderr":
		c.stderr.Write(data)
	}
	c.combined.Write(data)
	c.last = time.Now()
}

func (c *capture) snapshot() (string, string, string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.stdout.String(), c.stderr.String(), c.combined.String()
}

func (c *capture) lastActivity() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.last
}

func pump(wg *sync.WaitGroup, reader io.Reader, state *capture, stream string) {
	defer wg.Done()
	buf := make([]byte, 32*1024)
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			state.write(stream, buf[:n])
		}
		if err != nil {
			return
		}
	}
}

func waitProcess(ctx context.Context, cmd *exec.Cmd, waitCh <-chan error, state *capture, req execution.Request) (execution.Result, error) {
	ticker := time.NewTicker(25 * time.Millisecond)
	defer ticker.Stop()
	absolute := time.NewTimer(durationOrNever(req.Timeout))
	defer absolute.Stop()
	for {
		select {
		case err := <-waitCh:
			return processResult(cmd, "", false), normalizeExit(err)
		case <-ctx.Done():
			result := terminate(cmd, waitCh, "cancelled", req.TerminateGrace)
			return result, ctx.Err()
		case <-absolute.C:
			result := terminate(cmd, waitCh, "absolute_timeout", req.TerminateGrace)
			result.TimedOut = true
			return result, ErrTimeout
		case <-ticker.C:
			if req.IdleTimeout > 0 && time.Since(state.lastActivity()) > req.IdleTimeout {
				result := terminate(cmd, waitCh, "idle_timeout", req.TerminateGrace)
				result.TimedOut = true
				return result, ErrIdle
			}
		}
	}
}

func terminate(cmd *exec.Cmd, waitCh <-chan error, reason string, grace time.Duration) execution.Result {
	if grace <= 0 {
		grace = 250 * time.Millisecond
	}
	_ = signalProcessGroup(cmd, syscall.SIGTERM)
	timer := time.NewTimer(grace)
	defer timer.Stop()
	select {
	case <-waitCh:
	case <-timer.C:
		_ = signalProcessGroup(cmd, syscall.SIGKILL)
		<-waitCh
	}
	result := processResult(cmd, reason, true)
	return result
}

func signalProcessGroup(cmd *exec.Cmd, sig syscall.Signal) error {
	if cmd == nil || cmd.Process == nil {
		return nil
	}
	return syscall.Kill(-cmd.Process.Pid, sig)
}

func processResult(cmd *exec.Cmd, killReason string, timedOut bool) execution.Result {
	exit := 0
	if cmd == nil || cmd.ProcessState == nil {
		exit = -1
	} else {
		exit = cmd.ProcessState.ExitCode()
	}
	return execution.Result{ExitCode: exit, KillReason: killReason, TimedOut: timedOut}
}

func normalizeExit(err error) error {
	if err == nil {
		return nil
	}
	var exit *exec.ExitError
	if errors.As(err, &exit) {
		return nil
	}
	return err
}

func durationOrNever(duration time.Duration) time.Duration {
	if duration <= 0 {
		return 100 * 365 * 24 * time.Hour
	}
	return duration
}

func (r Runner) writeDiagnostic(command string, result execution.Result) (string, error) {
	if err := os.MkdirAll(r.DiagnosticsDir, 0o755); err != nil {
		return "", err
	}
	name := fmt.Sprintf("command-%d.txt", time.Now().UnixNano())
	path := filepath.Join(r.DiagnosticsDir, name)
	data := []byte(
		"command: " + command +
			"\nexit_code: " + fmt.Sprint(result.ExitCode) +
			"\nkill_reason: " + result.KillReason +
			"\n\nstdout:\n" + result.Stdout +
			"\n\nstderr:\n" + result.Stderr,
	)
	return path, os.WriteFile(path, data, 0o644)
}
