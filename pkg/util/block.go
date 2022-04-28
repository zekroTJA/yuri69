package util

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func Block(ctx context.Context) {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	select {
	case <-ctx.Done():
	case <-sc:
	}
	return
}
