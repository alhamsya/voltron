package meter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"

	modelPostgresql "github.com/alhamsya/voltron/internal/core/domain/postgresql"
	modelRequest "github.com/alhamsya/voltron/internal/core/domain/request"
	modelResponse "github.com/alhamsya/voltron/internal/core/domain/response"
)

func (s *Service) Reading(ctx context.Context, param []modelRequest.PowerMater) (modelResponse.Common, error) {
	return modelResponse.Common{}, errors.New("some error")
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
			DeviceID:  "iot",
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
