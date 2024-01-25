package postgresql

import (
	"dmaas/internal/app/dmaas/entity"
	"fmt"
	"gorm.io/gorm"
)

const FdwName = "postgres_fdw"

type PostgreSQLDriver struct {
	DB *gorm.DB
}

func (driver *PostgreSQLDriver) ImportDatabase(source entity.Source, localSchemaName string) error {
	serverName := driver.generateServerName(source)

	err := driver.createForeignServer(serverName, source)
	err = driver.createUserMapping(source)
	err = driver.importForeignSchema(localSchemaName, source)

	return err
}

func (driver *PostgreSQLDriver) createForeignServer(serverName string, source entity.Source) error {
	sql := fmt.Sprintf("CREATE SERVER IF NOT EXISTS %v "+
		"FOREIGN DATA WRAPPER %v "+
		"OPTIONS (host '%v', dbname '%v', port '%v');", serverName, FdwName, source.Host, source.Name, source.Port)

	return driver.DB.Exec(sql).Error
}

func (driver *PostgreSQLDriver) createUserMapping(source entity.Source) error {
	serverName := driver.generateServerName(source)

	sql := fmt.Sprintf("CREATE USER MAPPING IF NOT EXISTS "+
		"FOR CURRENT_USER SERVER %v OPTIONS (user '%v', password '%v')", serverName, source.Username, source.Password)
	return driver.DB.Exec(sql).Error
}

func (driver *PostgreSQLDriver) importForeignSchema(localSchema string, source entity.Source) error {
	serverName := driver.generateServerName(source)
	sql := fmt.Sprintf("IMPORT FOREIGN SCHEMA %v FROM SERVER %v INTO %v", source.Schema, serverName, localSchema)
	return driver.DB.Exec(sql).Error
}

func (driver *PostgreSQLDriver) generateServerName(source entity.Source) string {
	return fmt.Sprintf("postgresql_server_%v", source.ID)
}
