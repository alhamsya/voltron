package port

import (
	"context"
	modelRequest "github.com/alhamsya/voltron/internal/core/domain/request"
	"time"

	modelPostgresql "github.com/alhamsya/voltron/internal/core/domain/postgresql"
)

//go:generate mockgen -package=repomock -source=$GOFILE -destination=../../mock/repository/$GOFILE
type TimescaleRepo interface {
	BulkPowerMeter(ctx context.Context, param []modelPostgresql.PowerMeter) error
	GetMeterTimeSeries(ctx context.Context, param *modelRequest.TimeSeries) ([]modelPostgresql.PowerMeter, error)
	GetMeterLatestMeter(ctx context.Context, deviceID string) ([]modelPostgresql.Latest, error)
	GetMeterDailyUsage(ctx context.Context, deviceID string, from, to time.Time) ([]modelPostgresql.DailyUsage, error)
}
