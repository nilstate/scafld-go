package providers

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/nilstate/scafld-go/internal/core/execution"
	"github.com/nilstate/scafld-go/internal/core/review"
)

func TestProviderContract(t *testing.T) {
	t.Parallel()
	packet, err := (LocalProvider{Messages: []string{`{"type":"finding","severity":"blocking","summary":"bug"}`}}).Invoke(context.Background(), "task")
	if err != nil {
		t.Fatal(err)
	}
	if packet.Verdict != "fail" || len(packet.Findings) != 1 {
		t.Fatalf("packet = %+v", packet)
	}
}

type fakeRunner struct {
	result execution.Result
	err    error
	req    execution.Request
}

func (f *fakeRunner) Run(_ context.Context, req execution.Request) (execution.Result, error) {
	f.req = req
	return f.result, f.err
}

func TestCommandProviderParsesStdoutOnlyAndPassesTimeouts(t *testing.T) {
	t.Parallel()

	runner := &fakeRunner{result: execution.Result{
		ExitCode: 0,
		Stdout:   `{"type":"finding","severity":"blocking","summary":"bug"}` + "\n",
		Stderr:   "progress\n",
		Output:   "progress should not be parsed",
	}}
	packet, err := (CommandProvider{
		Command:     "reviewer",
		CWD:         "/tmp/work",
		Runner:      runner,
		Timeout:     time.Minute,
		IdleTimeout: time.Second,
	}).Invoke(context.Background(), "task")
	if err != nil {
		t.Fatal(err)
	}
	if packet.Verdict != review.VerdictFail || len(packet.Findings) != 1 {
		t.Fatalf("packet = %+v", packet)
	}
	if runner.req.CWD != "/tmp/work" || runner.req.Timeout != time.Minute || runner.req.IdleTimeout != time.Second {
		t.Fatalf("request = %+v", runner.req)
	}
}

func TestCommandProviderFailsClosedOnMissingRunnerInvalidOutputAndCleanPacketNonzeroExit(t *testing.T) {
	t.Parallel()

	if _, err := (CommandProvider{Command: "reviewer"}).Invoke(context.Background(), "task"); !errors.Is(err, ErrProviderFailed) {
		t.Fatalf("missing runner err = %v", err)
	}
	if _, err := (CommandProvider{Command: "reviewer", Runner: &fakeRunner{result: execution.Result{Stdout: "{invalid\n"}}}).Invoke(context.Background(), "task"); !errors.Is(err, review.ErrInvalidPacket) {
		t.Fatalf("invalid output err = %v", err)
	}
	runner := &fakeRunner{result: execution.Result{ExitCode: 1, Stdout: `{"type":"verdict","verdict":"pass"}` + "\n"}}
	if _, err := (CommandProvider{Command: "reviewer", Runner: runner}).Invoke(context.Background(), "task"); !errors.Is(err, ErrProviderFailed) {
		t.Fatalf("nonzero clean packet err = %v", err)
	}
}
