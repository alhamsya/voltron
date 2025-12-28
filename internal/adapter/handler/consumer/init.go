package consumer

import (
	port "github.com/alhamsya/voltron/internal/core/port/service"
)

type HandlerMeter struct {
	MeterService port.MeterService
}
