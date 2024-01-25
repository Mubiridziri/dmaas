package database

import (
	"dmaas/internal/app/dmaas/config"
	"dmaas/internal/app/dmaas/entity"
	"errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectAndMigrate() error {
	database, err := gorm.Open(postgres.Open(config.CFG.Database.DSN))
	if err != nil {
		return errors.New("failed connect to database")
	}

	err = database.AutoMigrate(entity.User{})
	err = database.AutoMigrate(entity.Field{})
	err = database.AutoMigrate(entity.Table{})
	err = database.AutoMigrate(entity.Source{})

	//DATA WRAPPERS
	err = database.Exec("CREATE EXTENSION IF NOT EXISTS postgres_fdw").Error

	if err != nil {
		return errors.New("failed auto migrate database")
	}

	DB = database

	return nil
}
