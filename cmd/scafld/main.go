package main

import (
	"os"

	"github.com/nilstate/scafld-go/internal/adapters/cli"
	"github.com/nilstate/scafld-go/internal/platform/signal"
)

func main() {
	ctx, handler := signal.RootContextWithOptions(nil, signal.Options{Escalate: func() { os.Exit(130) }})
	defer handler.Stop()
	os.Exit(cli.Run(ctx, os.Args[1:], os.Stdout, os.Stderr))
}
