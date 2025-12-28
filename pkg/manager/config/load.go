package config

import (
	"context"
	"os"
	"strconv"

	"github.com/alhamsya/voltron/internal/core/domain/config"
	"github.com/pkg/errors"
)

type Application struct {
	Credential config.Credential `mapstructure:"credential"`
	Static     config.Static     `mapstructure:"static"`
}

func GetConfig(ctx context.Context) *Application {
	cfg := &Application{
		Credential: config.Credential{
			ServiceSpecific: map[string]config.DBCredential{
				"timescale": {
					Primary: config.DBConnCredential{
						Username: os.Getenv("PGUSER"),
						Password: os.Getenv("PGPASSWORD"),
					},
					Replica: config.DBConnCredential{
						Username: os.Getenv("PGUSER"),
						Password: os.Getenv("PGPASSWORD"),
					},
				},
			},
		},
		Static: config.Static{
			Env: os.Getenv("ENV"),
			Frontend: &config.Frontend{
				URL: os.Getenv("FE_WEDDING_URL"),
			},
			App: &config.App{
				Rest: config.Rest{
					Port: convertStringToInt(os.Getenv("PORT")),
				},
			},
			ServiceSpecific: map[string]config.DBStatic{
				"timescale": {
					Primary: config.DBConnStatic{
						Host: os.Getenv("PGHOST"),
						Port: convertStringToInt(os.Getenv("PGPORT")),
						Name: os.Getenv("PGDATABASE"),
					},
					Replica: config.DBConnStatic{
						Host: os.Getenv("PGHOST"),
						Port: convertStringToInt(os.Getenv("PGPORT")),
						Name: os.Getenv("PGDATABASE"),
					},
				},
			},
		},
	}
	return cfg
}

func convertStringToInt(str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		panic(errors.New("failed strconv Atoi"))
	}

	return num
}
