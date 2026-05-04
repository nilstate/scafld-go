package git

import (
	"context"
	"os/exec"
	"strings"
)

type State struct {
	Changed []string
}

type Adapter struct {
	Root string
}

func (a Adapter) Status(ctx context.Context) (State, error) {
	cmd := exec.CommandContext(ctx, "git", "status", "--porcelain")
	cmd.Dir = a.Root
	out, err := cmd.Output()
	if err != nil {
		return State{}, err
	}
	var changed []string
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if len(line) > 3 {
			changed = append(changed, strings.TrimSpace(line[3:]))
		}
	}
	return State{Changed: changed}, nil
}

func (a Adapter) ChangedFiles(ctx context.Context) ([]string, error) {
	state, err := a.Status(ctx)
	if err != nil {
		return nil, nil
	}
	return state.Changed, nil
}
