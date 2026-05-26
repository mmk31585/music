package app

import (
	"fmt"
	"net/http"
)

func (a *App) NewHTTPServer() *http.Server {
	addr := fmt.Sprintf("%s:%s", a.Config.App.Host, a.Config.App.Port)

	return &http.Server{
		Addr:         addr,
		Handler:      a.Router,
		ReadTimeout:  a.Config.App.ReadTimeout,
		WriteTimeout: a.Config.App.WriteTimeout,
		IdleTimeout:  a.Config.App.IdleTimeout,
	}
}
