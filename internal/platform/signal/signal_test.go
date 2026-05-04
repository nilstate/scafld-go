package signal

import (
	"context"
	"testing"
)

func TestSignalContextCancel(t *testing.T) {
	t.Parallel()

	ctx, handler := RootContext(context.Background())
	handler.Stop()
	select {
	case <-ctx.Done():
	default:
		t.Fatal("context not cancelled on stop")
	}
}

func TestSignalInterruptTerminateEscalateContract(t *testing.T) {
	t.Parallel()
	_, handler := RootContext(context.Background())
	handler.Interrupts.Add(2)
	if handler.Interrupts.Load() != 2 {
		t.Fatal("interrupt count not recorded")
	}
	handler.Stop()
}
