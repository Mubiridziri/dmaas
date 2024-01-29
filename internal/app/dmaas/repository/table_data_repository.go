package repository

import (
	"dmaas/internal/app/dmaas/dto"
	"dmaas/internal/app/dmaas/entity"
	"fmt"
	"gorm.io/gorm"
)

type TableDataRepositoryInterface interface {
	ListTableData(localSchemaName string, table entity.Table, pagination dto.Query) ([]map[string]interface{}, error)
	GetCount(localSchemaName string, table entity.Table) int64
}

type TableDataRepository struct {
	DB *gorm.DB
}

func (repository *TableDataRepository) ListTableData(localSchemaName string, table entity.Table, pagination dto.Query) ([]map[string]interface{}, error) {
	offset := (pagination.Page - 1) * pagination.Limit
	var data []map[string]interface{}
	sql := fmt.Sprintf("SELECT * FROM %v.%v OFFSET %v LIMIT %v", localSchemaName, table.Name, offset, pagination.Limit)
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
