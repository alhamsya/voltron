package consumer

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"

	modelRequest "github.com/alhamsya/voltron/internal/core/domain/request"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (h *HandlerMeter) Consume(ctx context.Context, msg amqp.Delivery) error {
	var powerMeter []modelRequest.PowerMater
	if err := json.Unmarshal(msg.Body, &powerMeter); err != nil {
		return errors.Wrap(err, "failed json unmarshal")
	}

	err := h.MeterService.LogPowerMeter(ctx, powerMeter)
	if err != nil {
		return errors.Wrap(err, "failed service reading")
	}

	return nil
}
