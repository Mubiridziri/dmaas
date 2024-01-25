package config

import (
	"dmaas/internal/app/dmaas/utils"
	"errors"
	"fmt"
)

const DatabaseDsn = "DATABASE_DSN"

var CFG *Config

type Config struct {
	Database Database
}

type Database struct {
	DSN string
}

func New() *Config {
	return &Config{
		Database: Database{
			DSN: utils.GetEnvStr(DatabaseDsn, ""),
		},
	}
}

func LoadConfig() error {
	var cfg *Config

	cfg = New()
	if err := validate(cfg); err != nil {
		return err
	}

	CFG = cfg

	return nil
}

func validate(cfg *Config) error {
	if cfg.Database.DSN == "" {
		return createNotNullEnvError(DatabaseDsn)
	}

	return nil
}

func createNotNullEnvError(envName string) error {
	return errors.New(fmt.Sprintf("env variable %v cannot be null", envName))
}
