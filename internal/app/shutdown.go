package app

import (
	"context"

	"go.uber.org/zap"
)

func (a *App) Shutdown(ctx context.Context) error {
	a.Logger.Info("app shutdown: cancelling background contexts")

	// Step 1: Cancel all derived contexts — this signals every goroutine
	// that was started via AddBackground() to stop.
	if a.cancel != nil {
		a.cancel()
	}

	// Step 2: Stop the HTTP server (drains in-flight requests)
	if a.HTTPServer != nil {
		a.Logger.Info("app shutdown: draining HTTP server")
		_ = a.HTTPServer.Shutdown(ctx)
	}

	// Step 3: Shut down the event bus (stops outbox relay)
	if a.Events != nil {
		a.Logger.Info("app shutdown: stopping event bus relay")
		_ = a.Events.Shutdown(ctx)
	}

	// Step 4: Wait for all tracked background goroutines (WS hub, periodic
	// enrichment, cleanup workers, etc.) with a timeout.
	done := make(chan struct{})
	go func() {
		a.backgroundWg.Wait()
		close(done)
	}()

	select {
	case <-done:
		a.Logger.Info("app shutdown: all background goroutines finished")
	case <-ctx.Done():
		a.Logger.Warn("app shutdown: timed out waiting for background goroutines",
			zap.Error(ctx.Err()),
		)
	}

	// Step 5: Close infrastructure connections
	if a.DB != nil {
		a.DB.Close()
	}

	if a.Redis != nil {
		_ = a.Redis.Close()
	}

	if a.Logger != nil {
		_ = a.Logger.Sync()
	}

	a.Logger.Info("app shutdown complete")
	return nil
}
