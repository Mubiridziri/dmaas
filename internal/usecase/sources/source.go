package sources

import (
	"dmaas/internal/entity"
	"dmaas/internal/usecase/driver"
	"dmaas/internal/usecase/driver/postgresql"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

const (
	ImportAction = "import_action"
	RemoveAction = "remove_action"
)

type SourceView struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Host      string    `json:"host"`
	Port      int       `json:"port"`
	Username  string    `json:"-"`
	Password  string    `json:"-"`
	Schema    string    `json:"schema"`
	Alive     bool      `json:"alive"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SourceTableView struct {
	ID      int              `json:"id"`
	Name    string           `json:"name"`
	Comment string           `json:"comment"`
	Fields  []TableFieldView `json:"fields"`
}

type TableFieldView struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Comment  string `json:"comment"`
	Nullable bool   `json:"nullable"`
}

type CreateOrUpdateSourceView struct {
	Title    string `json:"title"  binding:"required"`
	Name     string `json:"name"  binding:"required"`
	Type     string `json:"type"  binding:"required"`
	Host     string `json:"host"  binding:"required"`
	Port     int    `json:"port"  binding:"required"`
	Schema   string `json:"schema"`
	Username string `json:"username"  binding:"required"`
	Password string `json:"password"  binding:"required"`
}

type Repository interface {
	CreateSource(source *entity.Source) error
	UpdateSource(source *entity.Source) error
	RemoveSource(source *entity.Source) error
	ListSources(page, limit int) ([]entity.Source, error)
	GetSourceById(id int) (entity.Source, error)
	GetSourcesCount() int64
}

type Controller struct {
	db *gorm.DB
	Repository
	SourceSender chan Job
}

type Job struct {
	Action string
	Source entity.Source
}

type PaginatedSourcesList struct {
	Total   int64        `json:"total"`
	Entries []SourceView `json:"entries"`
}

type InformationSchemaTable struct {
	TableName string
}

type InformationSchemaColumn struct {
	ColumnName string
	DataType   string
	IsNullable string
}

func NewController(repo Repository, db *gorm.DB, sourceSender chan Job) *Controller {
	return &Controller{Repository: repo, db: db, SourceSender: sourceSender}
}

func (c Controller) CreateSource(input CreateOrUpdateSourceView) (SourceView, error) {
	source := entity.Source{
		Title:    input.Title,
		Name:     input.Name,
		Type:     input.Type,
		Host:     input.Host,
		Port:     input.Port,
		Username: input.Username,
		Password: input.Password,
		Schema:   input.Schema,
	}

	err := c.Repository.CreateSource(&source)

	if err != nil {
		return SourceView{}, err
	}

	c.SourceSender <- Job{
		Action: ImportAction,
		Source: source,
	}

	return fromDBSource(&source), nil
}

//func (c Controller) UpdateSource(source CreateOrUpdateSourceView) (SourceView, error) {
//	return SourceView{}, nil
//}

func (c Controller) RemoveSource(id int) (SourceView, error) {
	source, err := c.Repository.GetSourceById(id)

	if err != nil {
		return SourceView{}, err
	}

	c.SourceSender <- Job{
		Action: RemoveAction,
		Source: source,
	}

	err = c.Repository.RemoveSource(&source)

	return fromDBSource(&source), err
}

func (c Controller) ListSources(page, limit int) (PaginatedSourcesList, error) {
	sources, err := c.Repository.ListSources(page, limit)

	if err != nil {
		return PaginatedSourcesList{}, err
	}

	var outputSources []SourceView

	for _, source := range sources {
		outputSources = append(outputSources, fromDBSource(&source))
	}

	return PaginatedSourcesList{
		Total:   c.Repository.GetSourcesCount(),
		Entries: outputSources,
	}, nil
}

func (c Controller) GetSourceById(id int) (SourceView, error) {
	source, err := c.Repository.GetSourceById(id)

	if err != nil {
		return SourceView{}, err
	}

	return fromDBSource(&source), nil
}

// ImportDatabase only start in goroutine!
func (c Controller) ImportDatabase(source entity.Source) {
	localSchemaName, err := c.createLocalSchema(source)
	databaseDriver, err := c.getDriverByType(source.Type)

	if err != nil {
		return
	}

	err = databaseDriver.ImportDatabase(source, localSchemaName)

	if err != nil {
		return
	}

	err = c.importStructure(source, localSchemaName)
	if err == nil {
		source.Alive = true
		c.db.Save(&source)
	}
}

func (c Controller) DeleteDatabase(source entity.Source) {
	databaseDriver, err := c.getDriverByType(source.Type)
	if err != nil {
		return
	}
	_ = databaseDriver.DropForeignServer(source)
	_ = c.dropLocalSchema(source)
}

func (c Controller) getLocalSchemaName(source entity.Source) string {
	return fmt.Sprintf("import_schema_%v", source.ID)
}

func (c Controller) importStructure(source entity.Source, localSchemaName string) error {
	tables, err := c.getTables(localSchemaName)

	for _, table := range tables {
		tableDB := entity.Table{Name: table.TableName, SourceID: source.ID}
		fields, err := c.getTableFields(tableDB.Name, localSchemaName)
		c.db.Create(&tableDB)

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

			c.db.Create(&fieldDB)
		}
	}

	return err
}

func (c Controller) createLocalSchema(source entity.Source) (string, error) {
	schemaName := c.getLocalSchemaName(source)
	sql := fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %v", schemaName)
	return schemaName, c.db.Exec(sql).Error
}

func (c Controller) dropLocalSchema(source entity.Source) error {
	schemaName := c.getLocalSchemaName(source)
	sql := fmt.Sprintf("DROP SCHEMA IF EXISTS %v", schemaName)
	return c.db.Exec(sql).Error
}

func (c Controller) getTables(localSchemaName string) ([]InformationSchemaTable, error) {
	var tables []InformationSchemaTable
	sql := fmt.Sprintf("SELECT table_name FROM information_schema.tables WHERE table_schema='%v'", localSchemaName)
	return tables, c.db.Raw(sql).Scan(&tables).Error
}

func (c Controller) getTableFields(tableName, localSchema string) ([]InformationSchemaColumn, error) {
	var fields []InformationSchemaColumn
	sql := fmt.Sprintf("SELECT column_name, data_type, is_nullable "+
		"FROM information_schema.columns WHERE table_schema='%v' AND table_name = '%v'", localSchema, tableName)
	return fields, c.db.Raw(sql).Scan(&fields).Error
}

func (c Controller) StartHandler() {
	go c.handleSourceEvents(c.SourceSender)
}

func (c Controller) handleSourceEvents(events chan Job) {
	for event := range events {
		switch event.Action {
		case ImportAction:
			c.ImportDatabase(event.Source)
			break
		case RemoveAction:
			c.DeleteDatabase(event.Source)
			break
		}
	}
}

// GetDriverByType Choice driver by source type (planed: mysql, oracle, db2, linter, mssql, etc.)
func (c Controller) getDriverByType(_type string) (driver.DriverInterface, error) {
	switch _type {
	case entity.PostgreSQLType:
		return &postgresql.PostgreSQLDriver{DB: c.db}, nil
		//case entity.MySQLType:
		//	return c.MySQLDriver.ImportDatabase(source, localSchemaName)
		//case entity.OracleType:
		//	return c.OracleDriver.ImportDatabase(source, localSchemaName)
	}
	return nil, errors.New("invalid driver type")
}

func fromDBSource(u *entity.Source) SourceView {
	source := SourceView{
		ID:        u.ID,
		Title:     u.Title,
		Name:      u.Name,
		Type:      u.Type,
		Host:      u.Host,
		Port:      u.Port,
		Username:  u.Username,
		Password:  u.Password,
		Schema:    u.Schema,
		Alive:     u.Alive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}

	return source
}
