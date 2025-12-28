package rabbitmq

import (
	"context"

	"github.com/alhamsya/voltron/pkg/manager/logging"

	amqp "github.com/rabbitmq/amqp091-go"
)

func New(ctx context.Context, queue, consumer string, funcConsumer func(context.Context, amqp.Delivery) error) <-chan amqp.Delivery {
	connection, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}

	channel, err := connection.Channel()
	if err != nil {
		panic(err)
	}

	delivery, err := channel.ConsumeWithContext(ctx, queue, consumer, true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	for message := range delivery {
		errConsumer := funcConsumer(ctx, message)
		if errConsumer != nil {
			logging.FromContext(ctx).Error().Err(errConsumer).Msg("failed consumer")
		}
	}

	return delivery
}
