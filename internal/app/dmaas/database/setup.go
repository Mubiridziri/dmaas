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
	//Sources
	err = database.AutoMigrate(entity.Field{})
	err = database.AutoMigrate(entity.Table{})
	err = database.AutoMigrate(entity.Source{})
	//Dictionaries
	err = database.AutoMigrate(entity.DictionaryValue{})
	err = database.AutoMigrate(entity.DictionaryRow{})
	err = database.AutoMigrate(entity.DictionaryField{})
	err = database.AutoMigrate(entity.Dictionary{})

	err = database.Create(&entity.User{Name: "admin", Username: "admin", Password: "admin"}).Error

	//DATA WRAPPERS
	err = database.Exec("CREATE EXTENSION IF NOT EXISTS postgres_fdw").Error

	if err != nil {
		return nil, errors.New("failed auto migrate database")
	}

	return database, nil
}
