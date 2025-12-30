package rabbitmq

import (
	"context"
	"fmt"
	"github.com/alhamsya/voltron/pkg/manager/config"
	"github.com/pkg/errors"
	"strings"

	"github.com/alhamsya/voltron/pkg/manager/logging"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	Username     string
	Password     string
	Host         string
	Port         int
	Queue        string
	ConsumerName string
}

type RabbitMQ struct {
	Cfg *config.Application

	Channel *amqp.Channel
}

func NewConsumer(ctx context.Context, cfg *Config, funcConsumer func(context.Context, amqp.Delivery) error) error {
	ctx = logging.ContextWithMetadata(ctx, nil)
	ctx = logging.InjectMetadata(ctx, map[string]any{
		"rabbitmq": map[string]any{
			"queue":         cfg.Queue,
			"consumer_name": cfg.ConsumerName,
		},
	})

	url := fmt.Sprintf("amqp://%s:%s@%s:%d/",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
	)
	connection, err := amqp.Dial(url)
	if err != nil {
		return errors.Wrap(err, "failed Dial")
	}
	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		return errors.Wrap(err, "failed Channel")
	}

	delivery, err := channel.ConsumeWithContext(ctx, cfg.Queue, cfg.ConsumerName, false, false, false, false, nil)
	if err != nil {
		return errors.Wrap(err, "failed ConsumeWithContext")
	}

	md := logging.MetadataFromContext(ctx)
	logging.FromContext(ctx).Info().
		Interface("log_info", md.ToMap()).
		Msg("accepting message")

	for {
		select {
		case <-ctx.Done():
			logging.FromContext(ctx).Warn().
				Interface("info", md.ToMap()).
				Msg("consumer context cancelled, shutting down")
			return nil
		case message, ok := <-delivery:
			if !ok {
				logging.FromContext(ctx).Warn().
					Interface("info", md.ToMap()).
					Msg("delivery channel closed")
				return nil
			}
			ctx = logging.InjectMetadata(ctx, map[string]any{
				"message_body": string(message.Body),
			})
			errConsumer := funcConsumer(ctx, message)
			if errConsumer != nil {
				if strings.Contains(errConsumer.Error(), "failed json unmarshal") {
					_ = message.Reject(false)
					logging.FromContext(ctx).Error().Err(errConsumer).
						Interface("info", md.ToMap()).
						Msg("failed json unmarshal")
					continue
				}

				if errNAck := message.Nack(false, true); errNAck != nil {
					logging.FromContext(ctx).Error().Err(errNAck).
						Interface("info", md.ToMap()).
						Msg("failed NAck message")
					return errNAck
				}

				logging.FromContext(ctx).Error().Err(errConsumer).
					Interface("info", md.ToMap()).
					Msg("failed consume")
				continue
			}

			if errAck := message.Ack(false); errAck != nil {
				logging.FromContext(ctx).Error().Err(errAck).
					Interface("info", md.ToMap()).
					Msg("failed NAck message")
				return errAck
			}
			logging.FromContext(ctx).Info().
				Interface("info", md.ToMap()).
				Msg("success consume")
		}
	}
}

func NewPublisher(cfg *config.Application) *RabbitMQ {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/",
		cfg.Credential.RabbitMQ.Username,
		cfg.Credential.RabbitMQ.Password,
		cfg.Static.RabbitMQ.Host,
		cfg.Static.RabbitMQ.Port,
	)
	connection, err := amqp.Dial(url)
	if err != nil {
		panic(errors.Wrap(err, "failed amqp Dial"))
	}

	channel, err := connection.Channel()
	if err != nil {
		panic(errors.Wrap(err, "failed connection Channel"))
	}

	return &RabbitMQ{
		Cfg:     cfg,
		Channel: channel,
	}
}
