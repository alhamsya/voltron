package logging

import (
	"context"
	"github.com/rs/zerolog"
	"os"
)

func FromContext(ctx context.Context) *zerolog.Logger {
	logging := zerolog.New(os.Stderr).With().Stack().Ctx(ctx).Timestamp().Logger()
	return &logging
}

func FromContextInfo(ctx context.Context) *zerolog.Event {
	logger := zerolog.New(os.Stderr).With().Stack().Ctx(ctx).Timestamp().Logger()
	logging := logger.Info()
	return logging
}
