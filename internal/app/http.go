package app

import (
	"net/http"
)

func (a *App) NewHTTPServer() *http.Server {
	return &http.Server{
		Addr:         a.Config.App.Host + ":" + a.Config.App.Port,
		Handler:      a.Router,
		ReadTimeout:  a.Config.App.ReadTimeout,
		WriteTimeout: a.Config.App.WriteTimeout,
		IdleTimeout:  a.Config.App.IdleTimeout,
	}
}
