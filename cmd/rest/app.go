package rest

import (
	"context"
	"github.com/alhamsya/voltron/internal/adapter/handler/rest"
	"github.com/alhamsya/voltron/internal/core/service/meter"
	"github.com/alhamsya/voltron/pkg/manager/protocol"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alhamsya/voltron/pkg/manager/config"
	"github.com/rs/zerolog"

	_ "go.uber.org/automaxprocs" // Automatically set GOMAXPROCS to match Linux container CPU quota.
)

func RunApp(ctx context.Context) error { //nolint:nolintlint,funlen
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, os.Interrupt)

	logger := zerolog.New(os.Stderr).With().Stack().Ctx(ctx).Timestamp().Logger()

	/* === GENERAL === */
	cfg := config.GetConfig(ctx)

	/* === DEPENDENCY INJECTION === */
	flightService := meter.NewService(&meter.Service{
		Cfg: cfg,
		Log: logger,
	})

	// Init router
	server := protocol.Rest(ctx, &protocol.RESTService{
		Cfg: cfg,
		Interactor: &rest.Interactor{
			MeterService: flightService,
		},
	})
	if err := server.Run(); err != nil {
		log.Fatalln("[Rest] service not running", err)
	}

	return nil
}
