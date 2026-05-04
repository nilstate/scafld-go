package providers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/nilstate/scafld-go/internal/core/execution"
	"github.com/nilstate/scafld-go/internal/core/review"
)

var ErrProviderFailed = errors.New("provider failed")

type Runner interface {
	Run(context.Context, execution.Request) (execution.Result, error)
}

type LocalProvider struct {
	Messages []string
}

func (p LocalProvider) Invoke(ctx context.Context, taskID string) (review.Packet, error) {
	var lines []string
	for _, msg := range p.Messages {
		if err := ctx.Err(); err != nil {
			return review.Packet{}, err
		}
		lines = append(lines, msg)
	}
	if len(lines) == 0 {
		lines = []string{`{"type":"verdict","verdict":"pass"}`}
	}
	return review.ParseNDJSON(strings.Join(lines, "\n") + "\n")
}

type CommandProvider struct {
	Command     string
	CWD         string
	Env         []string
	Runner      Runner
	Timeout     time.Duration
	IdleTimeout time.Duration
}

func (p CommandProvider) Invoke(ctx context.Context, taskID string) (review.Packet, error) {
	if p.Runner == nil {
		return review.Packet{}, fmt.Errorf("%w: runner is required", ErrProviderFailed)
	}
	if strings.TrimSpace(p.Command) == "" {
		return review.Packet{}, fmt.Errorf("%w: command is required", ErrProviderFailed)
	}
	env := append([]string(nil), p.Env...)
	env = append(env, "SCAFLD_TASK_ID="+taskID)
	result, err := p.Runner.Run(ctx, execution.Request{
		Command:     p.Command,
		CWD:         p.CWD,
		Env:         env,
		Timeout:     p.Timeout,
		IdleTimeout: p.IdleTimeout,
	})
	if err != nil && strings.TrimSpace(result.Stdout) == "" {
		return review.Packet{}, fmt.Errorf("%w: %v", ErrProviderFailed, err)
	}
	packet, parseErr := review.ParseNDJSON(result.Stdout)
	if parseErr != nil {
		if err != nil {
			return review.Packet{}, fmt.Errorf("%w: %v", ErrProviderFailed, err)
		}
		return review.Packet{}, parseErr
	}
	if err != nil {
		return review.Packet{}, fmt.Errorf("%w: %v", ErrProviderFailed, err)
	}
	if result.ExitCode != 0 && packet.Verdict != review.VerdictFail {
		return review.Packet{}, fmt.Errorf("%w: exit code %d", ErrProviderFailed, result.ExitCode)
	}
	return packet, nil
}
