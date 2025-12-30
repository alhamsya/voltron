package port

import "context"

//go:generate mockgen -package=repomock -source=$GOFILE -destination=../../mock/repository/$GOFILE
type RabbitMQRepo interface {
	PushPowerMeter(ctx context.Context, msgBody []byte) error
}
