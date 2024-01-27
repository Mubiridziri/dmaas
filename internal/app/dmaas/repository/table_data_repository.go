package repository

import (
	"dmaas/internal/app/dmaas/entity"
	"fmt"
	"gorm.io/gorm"
)

type TableDataRepositoryInterface interface {
	ListTableData(localSchemaName string, table entity.Table, page, limit int) ([]map[string]interface{}, error)
	GetCount(localSchemaName string, table entity.Table) int64
}

type TableDataRepository struct {
	DB *gorm.DB
}

func (repository *TableDataRepository) ListTableData(localSchemaName string, table entity.Table, page, limit int) ([]map[string]interface{}, error) {
	offset := (page - 1) * limit
	var data []map[string]interface{}
	sql := fmt.Sprintf("SELECT * FROM %v.%v OFFSET %v LIMIT %v", localSchemaName, table.Name, offset, limit)
	err := repository.DB.Raw(sql).Scan(&data).Error

	return data, err
}

func (repository *TableDataRepository) GetCount(localSchemaName string, table entity.Table) int64 {
	var data map[string]int64
	sql := fmt.Sprintf("SELECT COUNT(*) as count FROM %v.%v", localSchemaName, table.Name)
	err := repository.DB.Raw(sql).Scan(&data).Error

	if err != nil {
		return 0
	}

	return data["count"]
}
