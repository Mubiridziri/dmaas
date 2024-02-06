package tabledata

import (
	"dmaas/internal/entity"
	"fmt"
)

type PaginatedTableDataList struct {
	Total   int64                    `json:"total"`
	Entries []map[string]interface{} `json:"entries"`
}

type Repository interface {
	ListTableData(localSchemaName string, table entity.Table, page, limit int) ([]map[string]interface{}, error)
	GetTableDataCount(localSchemaName string, table entity.Table) int64
}

type Controller struct {
	Repository
}

func NewController(repo Repository) *Controller {
	return &Controller{Repository: repo}
}

func (c Controller) ListTableData(source entity.Source, table entity.Table, page, limit int) (PaginatedTableDataList, error) {
	localSchemaName := c.getLocalSchemaName(source)
	rows, err := c.Repository.ListTableData(localSchemaName, table, page, limit)

	if err != nil {
		return PaginatedTableDataList{}, err
	}

	return PaginatedTableDataList{
		Total:   c.Repository.GetTableDataCount(localSchemaName, table),
		Entries: rows,
	}, nil
}

// getLocalSchemaName BAD! ALREADY EXISTS IN sources.Controller
func (c Controller) getLocalSchemaName(source entity.Source) string {
	return fmt.Sprintf("import_schema_%v", source.ID)
}
