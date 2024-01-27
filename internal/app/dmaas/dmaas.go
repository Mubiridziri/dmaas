package dmaas

import (
	"dmaas/internal/app/dmaas/config"
	"dmaas/internal/app/dmaas/controller"
	"dmaas/internal/app/dmaas/database"
	"dmaas/internal/app/dmaas/repository"
	sources "dmaas/internal/app/dmaas/service"
	"net/http"
)

// Application  ...
type Application struct {
	BindAddr string
}

func (application *Application) Run() error {
	configLoader := config.ConfigLoader{}

	cfg, err := configLoader.LoadConfig()
	if err != nil {
		return err
	}

	db, err := database.ConnectAndMigrate(cfg)

	if err != nil {
		return err
	}

	router := controller.Router{
		//Repositories
		SourceRepository:    &repository.SourceRepository{DB: db},
		UserRepository:      &repository.UserRepository{DB: db},
		TableRepository:     &repository.TableRepository{DB: db},
		TableDataRepository: &repository.TableDataRepository{DB: db},
		//Services
		SourceManager: &sources.SourceManager{DB: db},
	}

	return http.ListenAndServe(application.BindAddr, router.NewRouter())
}

func New(bindAddr string) *Application {
	return &Application{BindAddr: bindAddr}
}
