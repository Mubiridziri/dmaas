package database

import (
	"dmaas/internal/app/dmaas/config"
	"dmaas/internal/app/dmaas/entity"
	"errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectAndMigrate(cfg *config.Config) (*gorm.DB, error) {
	database, err := gorm.Open(postgres.Open(cfg.Database.DSN))
	if err != nil {
		return nil, errors.New("failed connect to database")
	}

	err = database.AutoMigrate(entity.User{})
	err = database.AutoMigrate(entity.Field{})
	err = database.AutoMigrate(entity.Table{})
	err = database.AutoMigrate(entity.Source{})

	//DATA WRAPPERS
	err = database.Exec("CREATE EXTENSION IF NOT EXISTS postgres_fdw").Error

	if err != nil {
		return nil, errors.New("failed auto migrate database")
	}

	return database, nil
}
