package sources

import (
	"dmaas/internal/app/dmaas/database"
	"dmaas/internal/app/dmaas/entity"
	"dmaas/internal/app/dmaas/service/postgresql"
	"fmt"
)

type InformationSchemaTable struct {
	TableName string
}

type InformationSchemaColumn struct {
	ColumnName string
	DataType   string
	IsNullable string
}

func HandleSourceCreated(source entity.Source) {
	localSchemaName, err := createLocalSchema(source)

	//Choice driver by source type (planed: mysql, oracle, db2, linter, mssql, etc)
	switch source.Type {
	case entity.PostgreSQLType:
		err = postgresql.ImportPostgreSQLDatabase(source, localSchemaName)
		if err != nil {
			return
		}
	default:
		//return if unsupported type
		return
	}

	err = importStructure(source, localSchemaName)
	if err == nil {
		source.Alive = true
		database.DB.Save(&source)
	}
}

func importStructure(source entity.Source, localSchemaName string) error {
	tables, err := getTables(localSchemaName)

	for _, table := range tables {
		tableDB := entity.Table{Name: table.TableName, SourceID: source.ID}
		fields, err := getTableFields(tableDB.Name, localSchemaName)
		database.DB.Create(&tableDB)

		if err != nil {
			continue
		}

		for _, field := range fields {
			fieldDB := entity.Field{
				Name:     field.ColumnName,
				Type:     field.DataType,
				Nullable: field.IsNullable == "YES",
				TableID:  tableDB.ID,
			}

			database.DB.Create(&fieldDB)
		}
	}

	return err
}

func createLocalSchema(source entity.Source) (string, error) {
	schemaName := fmt.Sprintf("import_schema_%v", source.ID)
	sql := fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %v", schemaName)
	return schemaName, database.DB.Exec(sql).Error
}

func getTables(localSchemaName string) ([]InformationSchemaTable, error) {
	var tables []InformationSchemaTable
	sql := fmt.Sprintf("SELECT table_name FROM information_schema.tables WHERE table_schema='%v'", localSchemaName)
	return tables, database.DB.Raw(sql).Scan(&tables).Error
}

func getTableFields(tableName, localSchema string) ([]InformationSchemaColumn, error) {
	var fields []InformationSchemaColumn
	sql := fmt.Sprintf("SELECT column_name, data_type, is_nullable "+
		"FROM information_schema.columns WHERE table_schema='%v' AND table_name = '%v'", localSchema, tableName)
	return fields, database.DB.Raw(sql).Scan(&fields).Error
}
