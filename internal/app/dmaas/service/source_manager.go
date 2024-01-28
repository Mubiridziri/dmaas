package sources

import (
	"dmaas/internal/app/dmaas/entity"
	"dmaas/internal/app/dmaas/service/driver"
	"dmaas/internal/app/dmaas/service/driver/postgresql"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type SourceManager struct {
	DB *gorm.DB
}

type InformationSchemaTable struct {
	TableName string
}

type InformationSchemaColumn struct {
	ColumnName string
	DataType   string
	IsNullable string
}

// ImportDatabase only start in goroutine!
func (manager *SourceManager) ImportDatabase(source entity.Source) {
	localSchemaName, err := manager.createLocalSchema(source)
	databaseDriver, err := manager.getDriverByType(source.Type)

	if err != nil {
		return
	}

	err = databaseDriver.ImportDatabase(source, localSchemaName)

	if err != nil {
		return
	}

	err = manager.importStructure(source, localSchemaName)
	if err == nil {
		source.Alive = true
		manager.DB.Save(&source)
	}
}

func (manager *SourceManager) DeleteDatabase(source entity.Source) {
	databaseDriver, err := manager.getDriverByType(source.Type)
	if err != nil {
		return
	}
	_ = databaseDriver.DropForeignServer(source)
	_ = manager.dropLocalSchema(source)
}

func (manager *SourceManager) GetLocalSchemaName(source entity.Source) string {
	return fmt.Sprintf("import_schema_%v", source.ID)
}

func (manager *SourceManager) importStructure(source entity.Source, localSchemaName string) error {
	tables, err := manager.getTables(localSchemaName)

	for _, table := range tables {
		tableDB := entity.Table{Name: table.TableName, SourceID: source.ID}
		fields, err := manager.getTableFields(tableDB.Name, localSchemaName)
		manager.DB.Create(&tableDB)

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

			manager.DB.Create(&fieldDB)
		}
	}

	return err
}

func (manager *SourceManager) createLocalSchema(source entity.Source) (string, error) {
	schemaName := manager.GetLocalSchemaName(source)
	sql := fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %v", schemaName)
	return schemaName, manager.DB.Exec(sql).Error
}

func (manager *SourceManager) dropLocalSchema(source entity.Source) error {
	schemaName := manager.GetLocalSchemaName(source)
	sql := fmt.Sprintf("DROP SCHEMA IF EXISTS %v", schemaName)
	return manager.DB.Exec(sql).Error
}

func (manager *SourceManager) getTables(localSchemaName string) ([]InformationSchemaTable, error) {
	var tables []InformationSchemaTable
	sql := fmt.Sprintf("SELECT table_name FROM information_schema.tables WHERE table_schema='%v'", localSchemaName)
	return tables, manager.DB.Raw(sql).Scan(&tables).Error
}

func (manager *SourceManager) getTableFields(tableName, localSchema string) ([]InformationSchemaColumn, error) {
	var fields []InformationSchemaColumn
	sql := fmt.Sprintf("SELECT column_name, data_type, is_nullable "+
		"FROM information_schema.columns WHERE table_schema='%v' AND table_name = '%v'", localSchema, tableName)
	return fields, manager.DB.Raw(sql).Scan(&fields).Error
}

// GetDriverByType Choice driver by source type (planed: mysql, oracle, db2, linter, mssql, etc.)
func (manager *SourceManager) getDriverByType(_type string) (driver.DriverInterface, error) {
	switch _type {
	case entity.PostgreSQLType:
		return &postgresql.PostgreSQLDriver{DB: manager.DB}, nil
		//case entity.MySQLType:
		//	return manager.MySQLDriver.ImportDatabase(source, localSchemaName)
		//case entity.OracleType:
		//	return manager.OracleDriver.ImportDatabase(source, localSchemaName)
	}
	return nil, errors.New("invalid driver type")
}
