package rabbitmq

import (
	"context"

	"github.com/pkg/errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (r *RabbitMQ) PushPowerMeter(ctx context.Context, msgBody []byte) error {
	message := amqp.Publishing{
		Headers: amqp.Table{
			"sample": "value",
		},
		Body: msgBody,
	}
	err := r.Channel.PublishWithContext(ctx, "power-meter", "reading-key", false, false, message)
	if err != nil {
		return errors.Wrap(err, "failed PublishWithContext")
	}

	return nil
}
