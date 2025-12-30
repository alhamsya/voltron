package port

import (
	"context"
	modelRequest "github.com/alhamsya/voltron/internal/core/domain/request"

	modelPostgresql "github.com/alhamsya/voltron/internal/core/domain/postgresql"
)

type TimescaleRepo interface {
	BulkPowerMeter(ctx context.Context, param []modelPostgresql.PowerMeter) error
	GetTimeSeries(ctx context.Context, param *modelRequest.TimeSeries) ([]modelPostgresql.PowerMeter, error)
	GetLatestMeter(ctx context.Context, deviceID string) ([]modelPostgresql.Latest, error)
}
