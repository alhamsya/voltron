package meter

import (
	"github.com/alhamsya/voltron/internal/core/port/repository"
	"github.com/alhamsya/voltron/pkg/manager/config"
	"github.com/rs/zerolog"
)

type Service struct {
	Cfg *config.Application
	Log zerolog.Logger

	TimescaleRepo port.TimescaleRepo
}

func NewService(param *Service) *Service {
	return &Service{
		Cfg: param.Cfg,
		Log: param.Log,

		TimescaleRepo: param.TimescaleRepo,
	}
}
