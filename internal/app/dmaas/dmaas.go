package dmaas

import (
	"dmaas/internal/app/dmaas/config"
	"dmaas/internal/app/dmaas/controller"
	"dmaas/internal/app/dmaas/database"
	"net/http"
)

// Application  ...
type Application struct {
	BindAddr string
}

func (a *Application) Run() error {
	err := config.LoadConfig()
	if err != nil {
		return err
	}

	err = database.ConnectAndMigrate()

	if err != nil {
		return err
	}

	return http.ListenAndServe(a.BindAddr, controller.NewRouter())
}

func New(bindAddr string) *Application {
	return &Application{BindAddr: bindAddr}
}
