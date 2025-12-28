package meter

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	modelPostgresql "github.com/alhamsya/voltron/internal/core/domain/postgresql"
	modelRequest "github.com/alhamsya/voltron/internal/core/domain/request"
	modelResponse "github.com/alhamsya/voltron/internal/core/domain/response"
)

func (s *Service) Reading(ctx context.Context, param []modelRequest.PowerMater) (modelResponse.Common, error) {
	var reqDB []modelPostgresql.PowerMeter

	for _, data := range param {
		pmReading := modelPostgresql.PowerMeter{
			Time:      data.Date,
			DeviceID:  "iot",
			Metric:    data.Name,
			Value:     data.Data,
			EventHash: fmt.Sprintf("iot-%s", data.Name),
		}
		reqDB = append(reqDB, pmReading)
	}

	err := s.TimescaleRepo.BulkPowerMeter(ctx, reqDB)
	if err != nil {
		return modelResponse.Common{HttpCode: http.StatusInternalServerError}, errors.Wrap(err, "failed repo BulkPowerMeter")
	}

	return modelResponse.Common{
		HttpCode: http.StatusOK,
		Data:     nil,
		Metadata: nil,
	}, nil
}
