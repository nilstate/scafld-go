package review

import (
	"context"
	"time"

	"github.com/nilstate/scafld-go/internal/core/reconcile"
	"github.com/nilstate/scafld-go/internal/core/review"
	"github.com/nilstate/scafld-go/internal/core/session"
	"github.com/nilstate/scafld-go/internal/core/spec"
)

type SpecStore interface {
	Load(context.Context, string) (spec.Model, string, error)
	Save(context.Context, string, spec.Model) error
}

type SessionStore interface {
	Append(context.Context, string, session.Entry, string) (session.Session, error)
	Load(context.Context, string) (session.Session, error)
}

type Provider interface {
	Invoke(context.Context, string) (review.Packet, error)
}

type Clock interface{ Now() time.Time }

type Output struct {
	TaskID   string
	Verdict  string
	Findings []review.Finding
}

func Run(ctx context.Context, specs SpecStore, sessions SessionStore, provider Provider, clock Clock, taskID string) (Output, error) {
	model, path, err := specs.Load(ctx, taskID)
	if err != nil {
		return Output{}, err
	}
	packet, err := provider.Invoke(ctx, taskID)
	if err != nil {
		return Output{}, err
	}
	if err := review.ValidatePacket(packet); err != nil {
		return Output{}, err
	}
	now := clock.Now().UTC().Format(time.RFC3339)
	model.Status = spec.StatusReview
	model.Review.Status = "completed"
	model.Review.Verdict = packet.Verdict
	model.CurrentState.ReviewGate = packet.Verdict
	next, command := nextForVerdict(model.TaskID, packet.Verdict)
	model.CurrentState.Next = next
	model.CurrentState.AllowedFollowUp = command
	ledger, err := sessions.Append(ctx, model.TaskID, session.Entry{Type: "review", Status: packet.Verdict, Reason: "review gate " + packet.Verdict}, now)
	if err != nil {
		return Output{}, err
	}
	if loaded, loadErr := sessions.Load(ctx, model.TaskID); loadErr == nil {
		ledger = loaded
	}
	model = reconcile.FromSession(model, ledger)
	model.Status = spec.StatusReview
	model.Review.Status = "completed"
	model.Review.Verdict = packet.Verdict
	model.CurrentState.ReviewGate = packet.Verdict
	model.CurrentState.Next = next
	model.CurrentState.AllowedFollowUp = command
	if err := specs.Save(ctx, path, model); err != nil {
		return Output{}, err
	}
	return Output{TaskID: model.TaskID, Verdict: packet.Verdict, Findings: packet.Findings}, nil
}

func nextForVerdict(taskID string, verdict string) (string, string) {
	if verdict == "pass" {
		return "complete", "scafld complete " + taskID
	}
	return "repair", "scafld handoff " + taskID
}
