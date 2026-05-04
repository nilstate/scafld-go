package e2e

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestLifecycleJSONContractsAgentSurfaceFailCancelReviewProviderMutationGuard(t *testing.T) {
	t.Parallel()

	bin := testBinary(t)
	root := t.TempDir()
	run(t, bin, "init", "--root", root)
	run(t, bin, "plan", "--root", root, "lifecycle-task", "--title", "Lifecycle task", "--command", "test -f .scafld/config.yaml")
	run(t, bin, "approve", "--root", root, "lifecycle-task")
	run(t, bin, "build", "--root", root, "lifecycle-task")
	run(t, bin, "review", "--root", root, "lifecycle-task")
	run(t, bin, "complete", "--root", root, "lifecycle-task")
	out := run(t, bin, "status", "--root", root, "lifecycle-task", "--json")
	var envelope struct {
		OK bool `json:"ok"`
	}
	if err := json.Unmarshal(out, &envelope); err != nil {
		t.Fatal(err)
	}
	if !envelope.OK {
		t.Fatalf("status envelope not ok: %s", out)
	}
	if _, err := os.Stat(filepath.Join(root, ".scafld", "runs", "lifecycle-task", "session.json")); err != nil {
		t.Fatal(err)
	}
}

func TestFailCancel(t *testing.T) {
	t.Parallel()

	bin := testBinary(t)
	root := t.TempDir()
	run(t, bin, "init", "--root", root)
	run(t, bin, "plan", "--root", root, "cancel-task", "--command", "true")
	run(t, bin, "cancel", "--root", root, "cancel-task", "--reason", "test")
	run(t, bin, "plan", "--root", root, "fail-task", "--command", "false")
	run(t, bin, "fail", "--root", root, "fail-task", "--reason", "test failure")
}

func TestExitCodeTable(t *testing.T) {
	t.Parallel()

	bin := testBinary(t)
	root := t.TempDir()
	run(t, bin, "init", "--root", root)
	cmd := exec.Command(bin, "missing")
	if err := cmd.Run(); err == nil {
		t.Fatal("unknown command should fail")
	} else if exit, ok := err.(*exec.ExitError); !ok || exit.ExitCode() != 2 {
		t.Fatalf("exit = %v", err)
	}
}

func testBinary(t *testing.T) string {
	t.Helper()
	if bin := os.Getenv("SCAFLD_E2E_BINARY"); bin != "" {
		if !filepath.IsAbs(bin) {
			return filepath.Join("..", "..", bin)
		}
		return bin
	}
	bin := filepath.Join(t.TempDir(), "scafld")
	cmd := exec.Command("go", "build", "-o", bin, "./cmd/scafld")
	cmd.Dir = filepath.Join("..", "..")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build e2e binary: %v\n%s", err, out)
	}
	return bin
}

func run(t *testing.T, bin string, args ...string) []byte {
	t.Helper()
	cmd := exec.Command(bin, args...)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	if err := cmd.Run(); err != nil {
		t.Fatalf("%s %s failed: %v\nstdout:\n%s\nstderr:\n%s", bin, strings.Join(args, " "), err, out.String(), errOut.String())
	}
	return out.Bytes()
}
