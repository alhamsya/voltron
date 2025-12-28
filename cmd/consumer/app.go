package consumer

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func RunApp(ctx context.Context) error { //nolint:nolintlint,funlen
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, os.Interrupt)

	return nil
}
