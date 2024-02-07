package tables

import (
	"dmaas/internal/entity"
	"dmaas/internal/usecase/sources"
)

type PaginatedTablesView struct {
	Total   int64 `json:"total"`
	Entries []sources.SourceTableView
}

type Repository interface {
	ListTables(Source entity.Source, page, limit int) ([]entity.Table, error)
	GetTableById(id int) (entity.Table, error)
	GetTablesCount(source entity.Source) int64
}

type Controller struct {
	Repository
}

func NewController(repo Repository) *Controller {
	return &Controller{Repository: repo}
}

func (c Controller) ListTables(source entity.Source, page, limit int) (PaginatedTablesView, error) {
	tables, err := c.Repository.ListTables(source, page, limit)

	if err != nil {
		return PaginatedTablesView{}, err
	}

	var tableViews []sources.SourceTableView
	for _, table := range tables {
		tableViews = append(tableViews, fromDBTable(&table))
	}

	return PaginatedTablesView{
		Total:   c.Repository.GetTablesCount(source),
		Entries: tableViews,
	}, nil

}

func (c Controller) GetTableById(id int) (sources.SourceTableView, error) {
	table, err := c.Repository.GetTableById(id)

	if err != nil {
		return sources.SourceTableView{}, err
	}

	return fromDBTable(&table), nil
}

func fromDBTable(table *entity.Table) sources.SourceTableView {
	tableView := sources.SourceTableView{
		ID:      table.ID,
		Name:    table.Name,
		Comment: table.Comment,
	}

	var fields []sources.TableFieldView

	for _, field := range table.Fields {
		fieldView := sources.TableFieldView{
			ID:       field.ID,
			Name:     field.Name,
			Type:     field.Type,
			Comment:  field.Comment,
			Nullable: field.Nullable,
		}
		fields = append(fields, fieldView)
	}
	tableView.Fields = fields
	return tableView
}
