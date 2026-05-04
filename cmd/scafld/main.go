package main

import (
	"os"

	"github.com/nilstate/scafld-go/internal/adapters/cli"
	"github.com/nilstate/scafld-go/internal/platform/signal"
)

func main() {
	ctx, handler := signal.RootContext(nil)
	defer handler.Stop()
	os.Exit(cli.Run(ctx, os.Args[1:], os.Stdout, os.Stderr))
}
