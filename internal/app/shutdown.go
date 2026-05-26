package app

import (
	"context"
)

func (a *App) Shutdown(ctx context.Context) error {
	if a.HTTPServer != nil {
		_ = a.HTTPServer.Shutdown(ctx)
	}

	if a.DB != nil {
		a.DB.Close()
	}

	if a.Redis != nil {
		_ = a.Redis.Close()
	}

	if a.Logger != nil {
		_ = a.Logger.Sync()
	}

	return nil
}
