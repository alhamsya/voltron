package consumer

import (
	"context"
	"github.com/alhamsya/voltron/internal/adapter/handler/consumer"
	"github.com/alhamsya/voltron/internal/adapter/postgresql"
	"github.com/alhamsya/voltron/internal/adapter/rabbitmq"
	"github.com/alhamsya/voltron/internal/core/service/meter"
	"github.com/alhamsya/voltron/pkg/manager/config"
	"os"
	"os/signal"
	"syscall"
)

func RunApp(ctx context.Context) error { //nolint:nolintlint,funlen
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

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

	/* === DEPENDENCY INJECTION === */
	meterService := meter.NewService(&meter.Service{
		Cfg: cfg,

		TimescaleRepo: dbRepo,
	})

	/* === HANDLER === */
	handler := consumer.New(&consumer.HandlerMeter{
		MeterService: meterService,
	})

	err := rabbitmq.NewConsumer(ctx, &rabbitmq.Config{
		Username:     cfg.Credential.RabbitMQ.Username,
		Password:     cfg.Credential.RabbitMQ.Password,
		Host:         cfg.Static.RabbitMQ.Host,
		Port:         cfg.Static.RabbitMQ.Port,
		Queue:        "reading",
		ConsumerName: "consumer-reading",
	}, handler.Consume)

	if err != nil {
		panic(err)
	}

	return nil
}
