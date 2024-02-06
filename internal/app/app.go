package app

import (
	"dmaas/internal/config"
	"dmaas/internal/database"
	"dmaas/internal/entity"
	"dmaas/internal/server"
	"dmaas/internal/usecase/dictionaries"
	"dmaas/internal/usecase/dictionarydata"
	"dmaas/internal/usecase/sources"
	"dmaas/internal/usecase/tabledata"
	"dmaas/internal/usecase/tables"
	"dmaas/internal/usecase/users"
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

	repo := entity.NewRepository(db)
	var sourceJobs = make(chan sources.Job)

	userController := users.NewController(repo)
	sourceController := sources.NewController(repo, db, sourceJobs)
	tableController := tables.NewController(repo)
	tableDataController := tabledata.NewController(repo)
	dictController := dictionaries.NewController(repo)
	dictDataController := dictionarydata.NewController(repo)

	s := server.New(server.Config{
		UserController:           userController,
		SourceController:         sourceController,
		TableController:          tableController,
		TableDataController:      tableDataController,
		DictionaryController:     dictController,
		DictionaryDataController: dictDataController,
	})

	sourceController.StartHandler()

	return http.ListenAndServe(a.BindAddr, s.Router)
}

func New(bindAddr string) *Application {
	return &Application{BindAddr: bindAddr}
}
