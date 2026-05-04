package main

import (
	"os"

	"github.com/nilstate/scafld-go/internal/adapters/cli"
	"github.com/nilstate/scafld-go/internal/platform/signal"
)

func main() {
	ctx, handler := signal.RootContext(nil)
	code := cli.Run(ctx, os.Args[1:], os.Stdout, os.Stderr)
	handler.Stop()
	if handler.Escalated() {
		code = 130
	}
	os.Exit(code)
}
