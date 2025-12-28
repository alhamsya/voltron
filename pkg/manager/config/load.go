package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/alhamsya/voltron/internal/core/domain/config"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

type Application struct {
	Credential modelConfig.Credential `mapstructure:"credential"`
	Static     modelConfig.Static     `mapstructure:"static"`
}

func GetConfigENV() *Application {
	_ = godotenv.Load()

	cfg := &Application{
		Credential: modelConfig.Credential{
			ServiceSpecific: map[string]modelConfig.DBCredential{
				"timescale": {
					Primary: modelConfig.DBConnCredential{
						Username: os.Getenv("PGUSER"),
						Password: os.Getenv("PGPASSWORD"),
					},
					Replica: modelConfig.DBConnCredential{
						Username: os.Getenv("PGUSER"),
						Password: os.Getenv("PGPASSWORD"),
					},
				},
			},
		},
		Static: modelConfig.Static{
			Env: os.Getenv("ENV"),
			Frontend: &modelConfig.Frontend{
				URL: os.Getenv("FE_WEDDING_URL"),
			},
			App: &modelConfig.App{
				Rest: modelConfig.Rest{
					Port: convertStringToInt(os.Getenv("PORT")),
				},
			},
			ServiceSpecific: map[string]modelConfig.DBStatic{
				"timescale": {
					Primary: modelConfig.DBConnStatic{
						Host: os.Getenv("PGHOST"),
						Port: convertStringToInt(os.Getenv("PGPORT")),
						Name: os.Getenv("PGDATABASE"),
					},
					Replica: modelConfig.DBConnStatic{
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
		panic(errors.Wrap(err, fmt.Sprintf("failed strconv Atoi: %s", str)))
	}

	return num
}
