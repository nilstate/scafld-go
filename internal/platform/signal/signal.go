package signal

import (
	"context"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
)

type Handler struct {
	Interrupts atomic.Int32
	stop       func()
}

func RootContext(parent context.Context) (context.Context, *Handler) {
	if parent == nil {
		parent = context.Background()
	}
	ctx, cancel := context.WithCancel(parent)
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	handler := &Handler{
		stop: func() {
			signal.Stop(ch)
			close(ch)
			cancel()
		},
	}
	go func() {
		for range ch {
			handler.Interrupts.Add(1)
			cancel()
		}
	}()
	return ctx, handler
}

func (h *Handler) Stop() {
	if h == nil || h.stop == nil {
		return
	}
	h.stop()
	h.stop = nil
}
