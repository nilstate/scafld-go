package review

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	corereview "github.com/nilstate/scafld-go/internal/core/review"
	"github.com/nilstate/scafld-go/internal/core/session"
	"github.com/nilstate/scafld-go/internal/core/spec"
	"github.com/nilstate/scafld-go/internal/testkit/providerfake"
)

type fakeSpecs struct{ model spec.Model }

func (f *fakeSpecs) Load(context.Context, string) (spec.Model, string, error) {
	return f.model, "task.md", nil
}
func (f *fakeSpecs) Save(_ context.Context, _ string, model spec.Model) error {
	f.model = model
	return nil
}

type fakeSessions struct{ ledger session.Session }

func (f *fakeSessions) Append(_ context.Context, taskID string, entry session.Entry, now string) (session.Session, error) {
	if f.ledger.TaskID == "" {
		f.ledger = session.New(taskID, now)
	}
	f.ledger = f.ledger.WithEntry(entry)
	return f.ledger, nil
}

func (f *fakeSessions) Load(context.Context, string) (session.Session, error) { return f.ledger, nil }

type fakeProvider struct{ packet corereview.Packet }

func (f fakeProvider) Invoke(context.Context, corereview.Request) (corereview.Packet, error) {
	return f.packet, nil
}

type promptProvider struct {
	req    corereview.Request
	packet corereview.Packet
}

func (p *promptProvider) Invoke(_ context.Context, req corereview.Request) (corereview.Packet, error) {
	p.req = req
	return p.packet, nil
}

type fakeClock struct{}

func (fakeClock) Now() time.Time { return time.Date(2026, 5, 4, 0, 0, 0, 0, time.UTC) }

func TestProviderVerdictDrivesReviewState(t *testing.T) {
	t.Parallel()

	specs := &fakeSpecs{model: spec.Model{TaskID: "task", Title: "Task"}}
	sessions := &fakeSessions{}
	out, err := Run(context.Background(), specs, sessions, fakeProvider{packet: corereview.Packet{
		Verdict:  "fail",
		Findings: []corereview.Finding{{ID: "f1", Severity: corereview.SeverityBlocking, Summary: "bug"}},
	}}, fakeClock{}, "task")
	if err != nil {
		t.Fatal(err)
	}
	if out.Verdict != "fail" || len(out.Findings) != 1 {
		t.Fatalf("output = %+v", out)
	}
	if specs.model.CurrentState.AllowedFollowUp != "scafld handoff task" {
		t.Fatalf("next action = %q", specs.model.CurrentState.AllowedFollowUp)
	}
}

func TestProviderTimeoutMutationInvalidOutputPacketRepairFindingSignal(t *testing.T) {
	t.Parallel()

	specs := &fakeSpecs{model: spec.Model{TaskID: "task", Title: "Task"}}
	sessions := &fakeSessions{}
	out, err := Run(context.Background(), specs, sessions, providerfake.Provider{Mode: providerfake.ModeMutation}, fakeClock{}, "task")
	if err != nil {
		t.Fatal(err)
	}
	if out.Verdict != "fail" {
		t.Fatalf("mutation should fail review: %+v", out)
	}
	if len(out.Findings) != 1 || out.Findings[0].ID != "workspace_mutation" {
		t.Fatalf("mutation finding = %+v", out.Findings)
	}
}

func TestReviewRejectsInvalidDirectProviderPacket(t *testing.T) {
	t.Parallel()

	specs := &fakeSpecs{model: spec.Model{TaskID: "task", Title: "Task"}}
	_, err := Run(context.Background(), specs, &fakeSessions{}, fakeProvider{packet: corereview.Packet{Verdict: "maybe"}}, fakeClock{}, "task")
	if !errors.Is(err, corereview.ErrInvalidPacket) {
		t.Fatalf("invalid provider packet err = %v", err)
	}
}

func TestReviewPromptCarriesTaskContractToProvider(t *testing.T) {
	t.Parallel()

	provider := &promptProvider{packet: corereview.Packet{Verdict: corereview.VerdictPass}}
	specs := &fakeSpecs{model: spec.Model{TaskID: "task", Title: "Task", Summary: "Review this", Objectives: []string{"Keep evidence"}, Acceptance: spec.Acceptance{Criteria: []spec.Criterion{{ID: "ac1", Command: "go test ./...", ExpectedKind: "exit_code_zero"}}}}}
	_, err := Run(context.Background(), specs, &fakeSessions{}, provider, fakeClock{}, "task")
	if err != nil {
		t.Fatal(err)
	}
	if provider.req.TaskID != "task" || !strings.Contains(provider.req.Prompt, "Review this") || !strings.Contains(provider.req.Prompt, "ac1") {
		t.Fatalf("provider request = %+v", provider.req)
	}
}
