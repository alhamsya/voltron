package protocol

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/alhamsya/voltron/internal/adapter/handler/rest"
	"github.com/alhamsya/voltron/internal/core/domain/constant"
	"github.com/alhamsya/voltron/pkg/manager/config"
	"github.com/alhamsya/voltron/pkg/manager/graceful"
	"github.com/alhamsya/voltron/pkg/manager/logging"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog"
)

type RESTService struct {
	Cfg        *config.Application
	app        *fiber.App
	Interactor *rest.Interactor
}

func Rest(ctx context.Context, this *RESTService) *RESTService {
	app := fiber.New(fiber.Config{
		ReadTimeout: this.Cfg.Static.App.Rest.ReadTimeout,
		IdleTimeout: this.Cfg.Static.App.Rest.IdleTimeout,
	})

	app.Use(
		customLogger(ctx, this.Cfg),
		customCORS(this.Cfg),
		//customLimiter(this.Cfg),
		customRecover(),
	)

	routeHandler := &rest.Handler{
		Cfg:        this.Cfg,
		App:        app,
		Interactor: this.Interactor,
	}
	rest.New(routeHandler).Register()
	return &RESTService{
		Cfg:        this.Cfg,
		app:        app,
		Interactor: this.Interactor,
	}
}

func customLogger(ctx context.Context, cfg *config.Application) fiber.Handler {
	zerolog.TimestampFieldName = "timestamp"
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.ErrorStackMarshaler = logging.MarshalStack
	logger := zerolog.New(os.Stderr).With().
		Stack().
		Timestamp().
		Ctx(ctx).
		Logger()

	loggerConfig := fiberzerolog.New(fiberzerolog.Config{
		Logger:          &logger,
		FieldsSnakeCase: true,
		Messages:        []string{"server error", "client error", "success"},
	})

	return loggerConfig
}

func customCORS(config *config.Application) fiber.Handler {
	if config.Static.Frontend.URL == "" {
		return cors.New()
	}
	origin, err := url.Parse(config.Static.Frontend.URL)
	if err != nil {
		panic(err)
	}

	corsConfig := cors.Config{
		AllowOrigins:     fmt.Sprintf("%s://%s/", origin.Scheme, origin.Host),
		AllowMethods:     "OPTIONS,GET,POST,PUT,DELETE",
		AllowHeaders:     "Content-Type, Authorization",
		ExposeHeaders:    "Cross-Origin-Opener-Policy, Cross-Origin-Embedder-Policy",
		AllowCredentials: true,
	}

	return cors.New(corsConfig)
}

func customLimiter(config *config.Application) fiber.Handler {
	limiterConfig := limiter.New(limiter.Config{
		Max:        config.Static.App.Rest.Limiter.Max,        // max count of connections
		Expiration: config.Static.App.Rest.Limiter.Expiration, // expiration time of the limit
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusTooManyRequests)
		},
		SkipFailedRequests:     false,
		SkipSuccessfulRequests: false,
	})

	return limiterConfig
}

func customRecover() fiber.Handler {
	recoverConfig := recover.New(recover.Config{
		Next: func(c *fiber.Ctx) bool {
			err := c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
				"message": "please try again",
				"data":    nil,
			})
			if err != nil {
				panic(err)
			}
			return false
		},
		EnableStackTrace:  true,
		StackTraceHandler: nil,
	})

	return recoverConfig
}

func (r *RESTService) Run() error {
	return graceful.ServeRESTWithFiber(r.app, fmt.Sprintf(":%d", r.Cfg.Static.App.Rest.Port), constant.DefaultServerHTTPGraceTimeout)
}
