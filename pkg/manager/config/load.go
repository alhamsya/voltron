package config

import (
	"context"

	"github.com/alhamsya/voltron/internal/core/domain/config"
)

type Application struct {
	Credential config.Credential `mapstructure:"credential"`
	Static     config.Static     `mapstructure:"static"`
}

func GetConfig(ctx context.Context) *Application {
	var cfg Application
	return &cfg
}
