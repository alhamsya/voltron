package rest

import (
	"context"
	"github.com/alhamsya/voltron/internal/adapter/rabbitmq"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alhamsya/voltron/internal/adapter/handler/rest"
	"github.com/alhamsya/voltron/internal/adapter/postgresql"
	"github.com/alhamsya/voltron/internal/core/service/meter"
	"github.com/alhamsya/voltron/pkg/manager/config"
	"github.com/alhamsya/voltron/pkg/manager/protocol"
	_ "go.uber.org/automaxprocs" // Automatically set GOMAXPROCS to match Linux container CPU quota.
)

func RunApp(ctx context.Context) error { //nolint:nolintlint,funlen
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, os.Interrupt)

	/* === GENERAL === */
	cfg := config.GetConfigENV()

	/* === DATABASE === */
	dbPrimary := postgresql.Connect(ctx, &postgresql.Config{
		Username: cfg.Credential.ServiceSpecific["timescale"].Primary.Username,
		Password: cfg.Credential.ServiceSpecific["timescale"].Primary.Password,
		Host:     cfg.Static.ServiceSpecific["timescale"].Primary.Host,
		Port:     cfg.Static.ServiceSpecific["timescale"].Primary.Port,
		Name:     cfg.Static.ServiceSpecific["timescale"].Primary.Name,
	})
	dbReplica := postgresql.Connect(ctx, &postgresql.Config{
		Username: cfg.Credential.ServiceSpecific["timescale"].Primary.Username,
		Password: cfg.Credential.ServiceSpecific["timescale"].Primary.Password,
		Host:     cfg.Static.ServiceSpecific["timescale"].Primary.Host,
		Port:     cfg.Static.ServiceSpecific["timescale"].Primary.Port,
		Name:     cfg.Static.ServiceSpecific["timescale"].Primary.Name,
	})
	dbRepo := postgresql.New(cfg, dbPrimary, dbReplica)

	rabbitMQRepo := rabbitmq.NewPublisher(cfg)

	/* === DEPENDENCY INJECTION === */
	meterService := meter.NewService(&meter.Service{
		Cfg: cfg,

		TimescaleRepo: dbRepo,
		RabbitMQRepo:  rabbitMQRepo,
	})

	/* === HANDLER === */
	server := protocol.Rest(ctx, &protocol.RESTService{
		Cfg: cfg,
		Interactor: &rest.Interactor{
			MeterService: meterService,
		},
	})
	if err := server.Run(); err != nil {
		log.Fatalln("[Rest] service not running", err)
	}

	return nil
}
