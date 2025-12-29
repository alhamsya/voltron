package meter

import (
	"github.com/alhamsya/voltron/internal/core/port/repository"
	"github.com/alhamsya/voltron/pkg/manager/config"
)

type Service struct {
	Cfg *config.Application

	TimescaleRepo port.TimescaleRepo
	RabbitMQRepo  port.RabbitMQRepo
}

func NewService(param *Service) *Service {
	return &Service{
		Cfg: param.Cfg,

		TimescaleRepo: param.TimescaleRepo,
		RabbitMQRepo:  param.RabbitMQRepo,
	}
}
