package port

import (
	"context"

	modelRequest "github.com/alhamsya/voltron/internal/core/domain/request"
	modelResponse "github.com/alhamsya/voltron/internal/core/domain/response"
)

type MeterService interface {
	Reading(ctx context.Context, param []modelRequest.PowerMater) (modelResponse.Common, error)
}
