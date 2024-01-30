package usecase

import (
	"dmaas/internal/app/dmaas/dto"
	"dmaas/internal/app/dmaas/entity"
	"dmaas/internal/app/dmaas/repository"
	"dmaas/internal/app/dmaas/usecase/driver"
	"dmaas/internal/app/dmaas/usecase/driver/postgresql"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type SourceUseCaseInterface interface {
	CreateSource(source *entity.Source) error
	UpdateSource(source *entity.Source) error
	RemoveSource(source *entity.Source) error
	ListSources(pagination dto.Query) ([]entity.Source, error)
	GetSourceById(id int) (entity.Source, error)
	GetCount() int64
	ImportDatabase(source entity.Source)
	DeleteDatabase(source entity.Source)
}

type SourceUseCase struct {
	DB               *gorm.DB
	SourceRepository repository.SourceRepositoryInterface
	SourceSender     chan dto.SourceMessage
}

type InformationSchemaTable struct {
	TableName string
}

type InformationSchemaColumn struct {
	ColumnName string
	DataType   string
	IsNullable string
}

func (useCase *SourceUseCase) CreateSource(source *entity.Source) error {
	err := useCase.SourceRepository.CreateSource(source)
	//go useCase.importDatabase(*source)

	useCase.SourceSender <- dto.SourceMessage{
		Source: source,
		Action: dto.ImportDatabaseAction,
	}

	return err
}

func (useCase *SourceUseCase) UpdateSource(source *entity.Source) error {
	return useCase.SourceRepository.UpdateSource(source)
}

func (useCase *SourceUseCase) RemoveSource(source *entity.Source) error {
	err := useCase.SourceRepository.RemoveSource(source)

	useCase.SourceSender <- dto.SourceMessage{
		Source: source,
		Action: dto.RemoveDatabaseAction,
	}

	return err
}

func (useCase *SourceUseCase) ListSources(pagination dto.Query) ([]entity.Source, error) {
	return useCase.SourceRepository.ListSources(pagination)
}

func (useCase *SourceUseCase) GetSourceById(id int) (entity.Source, error) {
	return useCase.SourceRepository.GetSourceById(id)
}

func (useCase *SourceUseCase) GetCount() int64 {
	return useCase.SourceRepository.GetCount()
}

// ImportDatabase only start in goroutine!
func (useCase *SourceUseCase) ImportDatabase(source entity.Source) {
	localSchemaName, err := useCase.createLocalSchema(source)
	databaseDriver, err := useCase.getDriverByType(source.Type)

	if err != nil {
		return
	}

	err = databaseDriver.ImportDatabase(source, localSchemaName)

	if err != nil {
		return
	}

	err = useCase.importStructure(source, localSchemaName)
	if err == nil {
		source.Alive = true
		useCase.DB.Save(&source)
	}
}

func (useCase *SourceUseCase) DeleteDatabase(source entity.Source) {
	databaseDriver, err := useCase.getDriverByType(source.Type)
	if err != nil {
		return
	}
	_ = databaseDriver.DropForeignServer(source)
	_ = useCase.dropLocalSchema(source)
}

func (useCase *SourceUseCase) getLocalSchemaName(source entity.Source) string {
	return fmt.Sprintf("import_schema_%v", source.ID)
}

func (useCase *SourceUseCase) importStructure(source entity.Source, localSchemaName string) error {
	tables, err := useCase.getTables(localSchemaName)

	for _, table := range tables {
		tableDB := entity.Table{Name: table.TableName, SourceID: source.ID}
		fields, err := useCase.getTableFields(tableDB.Name, localSchemaName)
		useCase.DB.Create(&tableDB)

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

			useCase.DB.Create(&fieldDB)
		}
	}

	return err
}

func (useCase *SourceUseCase) createLocalSchema(source entity.Source) (string, error) {
	schemaName := useCase.getLocalSchemaName(source)
	sql := fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %v", schemaName)
	return schemaName, useCase.DB.Exec(sql).Error
}

func (useCase *SourceUseCase) dropLocalSchema(source entity.Source) error {
	schemaName := useCase.getLocalSchemaName(source)
	sql := fmt.Sprintf("DROP SCHEMA IF EXISTS %v", schemaName)
	return useCase.DB.Exec(sql).Error
}

func (useCase *SourceUseCase) getTables(localSchemaName string) ([]InformationSchemaTable, error) {
	var tables []InformationSchemaTable
	sql := fmt.Sprintf("SELECT table_name FROM information_schema.tables WHERE table_schema='%v'", localSchemaName)
	return tables, useCase.DB.Raw(sql).Scan(&tables).Error
}

func (useCase *SourceUseCase) getTableFields(tableName, localSchema string) ([]InformationSchemaColumn, error) {
	var fields []InformationSchemaColumn
	sql := fmt.Sprintf("SELECT column_name, data_type, is_nullable "+
		"FROM information_schema.columns WHERE table_schema='%v' AND table_name = '%v'", localSchema, tableName)
	return fields, useCase.DB.Raw(sql).Scan(&fields).Error
}

// GetDriverByType Choice driver by source type (planed: mysql, oracle, db2, linter, mssql, etc.)
func (useCase *SourceUseCase) getDriverByType(_type string) (driver.DriverInterface, error) {
	switch _type {
	case entity.PostgreSQLType:
		return &postgresql.PostgreSQLDriver{DB: useCase.DB}, nil
		//case entity.MySQLType:
		//	return useCase.MySQLDriver.ImportDatabase(source, localSchemaName)
		//case entity.OracleType:
		//	return useCase.OracleDriver.ImportDatabase(source, localSchemaName)
	}
	return nil, errors.New("invalid driver type")
}
