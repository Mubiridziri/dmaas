package app

import (
	"dmaas/internal/config"
	"dmaas/internal/context"
	"dmaas/internal/controller"
	"dmaas/internal/database"
	"dmaas/internal/dto"
	"dmaas/internal/handler"
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
