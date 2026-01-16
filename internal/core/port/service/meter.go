package port

import (
	"context"
	modelPower "github.com/alhamsya/voltron/internal/core/domain/power"
	"time"

	modelRequest "github.com/alhamsya/voltron/internal/core/domain/request"
	modelResponse "github.com/alhamsya/voltron/internal/core/domain/response"
)

//go:generate mockgen -package=usecasemock -source=$GOFILE -destination=../../../mock/usecase/$GOFILE
type MeterService interface {
	Reading(ctx context.Context, param []modelRequest.PowerMater) (modelResponse.Common, error)
	LogPowerMeter(ctx context.Context, param []modelRequest.PowerMater) error
	TimeSeries(ctx context.Context, param *modelRequest.TimeSeries) (modelResponse.Common, error)
	Latest(ctx context.Context, deviceID string) (modelResponse.Common, error)
	DailyUsage(ctx context.Context, deviceID string, from, to time.Time) (modelResponse.Common, error)
	BuildBilling(ctx context.Context, deviceID string, from, to time.Time) (modelPower.BillingSummary, []modelPower.DailyLine, error)
}
