package port

import "context"

type RabbitMQRepo interface {
	PushPowerMeter(ctx context.Context, msgBody []byte) error
}
