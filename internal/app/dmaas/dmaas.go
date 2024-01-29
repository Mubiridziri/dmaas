package dmaas

import (
	"dmaas/internal/app/dmaas/config"
	"dmaas/internal/app/dmaas/context"
	"dmaas/internal/app/dmaas/controller"
	"dmaas/internal/app/dmaas/database"
	"dmaas/internal/app/dmaas/dto"
	"net/http"
)

// Application  ...
type Application struct {
	BindAddr string
}

func (a *Application) Run() error {
	configLoader := config.ConfigLoader{}

	cfg, err := configLoader.LoadConfig()
	if err != nil {
		return err
	}

	db, err := database.ConnectAndMigrate(cfg)

	if err != nil {
		return err
	}

	var sourceChan = make(chan dto.SourceChan)

	ctx := context.New(cfg, db, &sourceChan)
	router := controller.Router{Context: ctx}

	return http.ListenAndServe(a.BindAddr, router.NewRouter())
}

func New(bindAddr string) *Application {
	return &Application{BindAddr: bindAddr}
}
