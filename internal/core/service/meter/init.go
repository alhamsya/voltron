package meter

import (
	"github.com/alhamsya/voltron/pkg/manager/config"
	"github.com/rs/zerolog"
)

type Service struct {
	Cfg *config.Application
	Log zerolog.Logger
}

func NewService(param *Service) *Service {
	return &Service{
		Cfg: param.Cfg,
		Log: param.Log,
	}
}
