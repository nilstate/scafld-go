package providers

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestProviderContract(t *testing.T) {
	t.Parallel()
	var out bytes.Buffer
	err := (LocalProvider{Messages: []string{"one", "two"}}).Invoke(context.Background(), &out)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "two") {
		t.Fatalf("output %q", out.String())
	}
}
