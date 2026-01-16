package meter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/alhamsya/voltron/internal/core/domain/constant"
	"github.com/pkg/errors"

	modelPostgresql "github.com/alhamsya/voltron/internal/core/domain/postgresql"
	modelPower "github.com/alhamsya/voltron/internal/core/domain/power"
	modelRequest "github.com/alhamsya/voltron/internal/core/domain/request"
	modelResponse "github.com/alhamsya/voltron/internal/core/domain/response"
)

func (s *Service) Reading(ctx context.Context, param []modelRequest.PowerMater) (modelResponse.Common, error) {
	msgByte, err := json.Marshal(param)
	if err != nil {
		return modelResponse.Common{HttpCode: http.StatusInternalServerError}, errors.Wrap(err, "failed json marshal")
	}

	err = s.RabbitMQRepo.PushPowerMeter(ctx, msgByte)
	if err != nil {
		return modelResponse.Common{HttpCode: http.StatusInternalServerError}, errors.Wrap(err, "failed publish to RabbitMQ")
	}

	return modelResponse.Common{
		HttpCode: http.StatusOK,
		Data:     nil,
		Metadata: nil,
	}, nil
}

func (s *Service) LogPowerMeter(ctx context.Context, param []modelRequest.PowerMater) error {
	var reqDB []modelPostgresql.PowerMeter
	for _, data := range param {
		pmReading := modelPostgresql.PowerMeter{
			Time:      data.Date,
			DeviceID:  constant.DeviceIoT,
			Metric:    data.Name,
			Value:     data.Data,
			EventHash: fmt.Sprintf("iot-%s-%d", data.Name, time.Now().Unix()),
		}
		reqDB = append(reqDB, pmReading)
	}

	err := s.TimescaleRepo.BulkPowerMeter(ctx, reqDB)
	if err != nil {
		return errors.Wrap(err, "failed repo BulkPowerMeter")
	}

	return nil
}

func (s *Service) TimeSeries(ctx context.Context, param *modelRequest.TimeSeries) (modelResponse.Common, error) {
	respData, err := s.TimescaleRepo.GetMeterTimeSeries(ctx, param)
	if err != nil {
		return modelResponse.Common{}, errors.Wrap(err, "failed repo GetTimeSeries")
	}

	if respData == nil {
		return modelResponse.Common{
			HttpCode: http.StatusOK,
			Data:     []modelPostgresql.PowerMeter{},
		}, err
	}

	return modelResponse.Common{
		HttpCode: http.StatusOK,
		Data:     respData,
	}, nil
}

func (s *Service) Latest(ctx context.Context, deviceID string) (modelResponse.Common, error) {
	respData, err := s.TimescaleRepo.GetMeterLatestMeter(ctx, deviceID)
	if err != nil {
		return modelResponse.Common{}, errors.Wrap(err, "failed repo GetLatestMeter")
	}

	if respData == nil {
		return modelResponse.Common{
			HttpCode: http.StatusOK,
			Data:     []modelPostgresql.Latest{},
		}, err
	}

	return modelResponse.Common{
		HttpCode: http.StatusOK,
		Data:     respData,
	}, nil
}

func (s *Service) DailyUsage(ctx context.Context, deviceID string, from, to time.Time) (modelResponse.Common, error) {
	respData, err := s.TimescaleRepo.GetMeterDailyUsage(ctx, deviceID, from, to)
	if err != nil {
		return modelResponse.Common{}, errors.Wrap(err, "failed repo GetMeterDailyUsage")
	}

	if respData == nil {
		return modelResponse.Common{
			HttpCode: http.StatusOK,
			Data:     []modelPostgresql.DailyUsage{},
		}, err
	}

	return modelResponse.Common{
		HttpCode: http.StatusOK,
		Data:     respData,
	}, nil
}

// BuildBilling builds a billing summary for a device within a given time range.
// It calculates total energy usage (kWh), applies the configured rate and tax,
// and returns both the billing summary and daily usage line details.
func (s *Service) BuildBilling(ctx context.Context, deviceID string, from, to time.Time) (modelPower.BillingSummary, []modelPower.DailyLine, error) {
	// Retrieve total energy consumption (kWh) for the given device and period
	totalKwh, err := s.TimescaleRepo.GetBillingSummary(ctx, deviceID, from, to)
	if err != nil {
		return modelPower.BillingSummary{}, nil, errors.Wrap(err, "failed repo GetBillingSummary")
	}

	// Retrieve daily usage breakdown used for billing detail lines
	lines, err := s.TimescaleRepo.GetDailyUsageLines(ctx, deviceID, from, to)
	if err != nil {
		return modelPower.BillingSummary{}, nil, errors.Wrap(err, "failed repo GetDailyUsageLines")
	}

	// Calculate subtotal based on total kWh and configured rate
	subtotal := totalKwh * s.Cfg.Static.App.PowerMeter.Rate
	// Calculate tax from subtotal using configured tax rate
	tax := subtotal * s.Cfg.Static.App.PowerMeter.TaxRate
	// Calculate final total including tax
	total := subtotal + tax

	// Build billing summary response object
	sum := modelPower.BillingSummary{
		DeviceID:   deviceID,
		From:       from.Format(constant.DateOnly),
		To:         to.Format(constant.DateOnly),
		TotalKwh:   totalKwh,
		RatePerKwh: s.Cfg.Static.App.PowerMeter.Rate,
		Subtotal:   subtotal,
		Tax:        tax,
		Total:      total,
	}
	return sum, lines, nil
}
