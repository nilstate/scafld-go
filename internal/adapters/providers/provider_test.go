package providers

import (
	"context"
	"testing"
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
