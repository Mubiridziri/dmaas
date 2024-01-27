package repository

import (
	"dmaas/internal/app/dmaas/entity"
	"gorm.io/gorm"
)

type TableRepositoryInterface interface {
	ListTables(sourceId int, page, limit int) ([]entity.Table, error)
	GetTableById(id int) (entity.Table, error)
}

type TableRepository struct {
	DB *gorm.DB
}

func (repository *TableRepository) ListTables(sourceId int, page, limit int) ([]entity.Table, error) {
	var tables []entity.Table
	offset := (page - 1) * limit
	if err := repository.DB.
		Offset(offset).
		Limit(limit).
		Where(entity.Table{SourceID: sourceId}).
		Find(&tables).Error; err != nil {
		return []entity.Table{}, err
	}

	return tables, nil
}

func (repository *TableRepository) GetTableById(id int) (entity.Table, error) {
	var table entity.Table
	if err := repository.DB.
		Preload("Fields").
		Where(entity.Table{ID: id}).
		First(&table).Error; err != nil {
		return entity.Table{}, err
	}

	return table, nil
}
