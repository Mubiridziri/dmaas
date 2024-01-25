package postgresql

import (
	"dmaas/internal/app/dmaas/database"
	"dmaas/internal/app/dmaas/entity"
	"fmt"
)

const FdwName = "postgres_fdw"

func ImportPostgreSQLDatabase(source entity.Source, localSchemaName string) error {
	serverName := generateServerName(source)

	err := createForeignServer(serverName, source)
	err = createUserMapping(source)
	err = importForeignSchema(localSchemaName, source)

	return err
}

func createForeignServer(serverName string, source entity.Source) error {
	sql := fmt.Sprintf("CREATE SERVER IF NOT EXISTS %v "+
		"FOREIGN DATA WRAPPER %v "+
		"OPTIONS (host '%v', dbname '%v', port '%v');", serverName, FdwName, source.Host, source.Name, source.Port)

	return database.DB.Exec(sql).Error
}

func createUserMapping(source entity.Source) error {
	serverName := generateServerName(source)

	sql := fmt.Sprintf("CREATE USER MAPPING IF NOT EXISTS "+
		"FOR CURRENT_USER SERVER %v OPTIONS (user '%v', password '%v')", serverName, source.Username, source.Password)
	return database.DB.Exec(sql).Error
}

func importForeignSchema(localSchema string, source entity.Source) error {
	serverName := generateServerName(source)
	sql := fmt.Sprintf("IMPORT FOREIGN SCHEMA %v FROM SERVER %v INTO %v", source.Schema, serverName, localSchema)
	return database.DB.Exec(sql).Error
}

func generateServerName(source entity.Source) string {
	return fmt.Sprintf("postgresql_server_%v", source.ID)
}
