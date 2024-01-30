package dmaas

import (
	"dmaas/internal/app/dmaas/config"
	"dmaas/internal/app/dmaas/context"
	"dmaas/internal/app/dmaas/controller"
	"dmaas/internal/app/dmaas/database"
	"dmaas/internal/app/dmaas/dto"
	"dmaas/internal/app/dmaas/handler"
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

	var sourceSender = make(chan dto.SourceMessage)

	//Create ApplicationContext
	ctx := context.New(cfg, db, sourceSender)

	//Start SourceHandler
	sourceHandler := handler.SourceHandler{Context: ctx}
	go sourceHandler.HandleSources(sourceSender)

	router := controller.Router{Context: ctx}

	return http.ListenAndServe(a.BindAddr, router.NewRouter())
}

func New(bindAddr string) *Application {
	return &Application{BindAddr: bindAddr}
}
