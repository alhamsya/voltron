package port

import "context"

type TimescaleRepo interface {
	Get(ctx context.Context) error
}
