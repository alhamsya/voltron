package graceful

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

// ServeRESTWithFiber start the http server using fiber on the given address and add graceful shutdown handler
func ServeRESTWithFiber(app *fiber.App, address string, graceTimeout time.Duration) error {
	stoppedCh := WaitTermSig(func(ctx context.Context) error {
		if graceTimeout == 0 {
			graceTimeout = 10 * time.Second
		}
		stopped := make(chan struct{})
		ctx, cancel := context.WithTimeout(ctx, graceTimeout)
		defer cancel()
		go func() {
			if err := app.Shutdown(); err != nil {
				log.Printf("[Rest] app shutdown error: %v", err)
			}
			close(stopped)
		}()

		select {
		case <-ctx.Done():
			return errors.New("[Rest] server shutdown timed out")
		case <-stopped:

		}

		return nil
	})

	log.Printf("[Rest] server running on address: %v", address)

	// start serving
	if err := app.Listen(address); err != nil {
		return fmt.Errorf("net listen: %w", err)
	}

	<-stoppedCh
	log.Println("[Rest] server stopped")

	return nil
}
