package port

import (
	"context"

	modelPostgresql "github.com/alhamsya/voltron/internal/core/domain/postgresql"
)

type TimescaleRepo interface {
	BulkPowerMeter(ctx context.Context, param []modelPostgresql.PowerMeter) error
}
