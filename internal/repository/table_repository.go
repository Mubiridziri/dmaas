package repository

import (
	"dmaas/internal/dto"
	"dmaas/internal/entity"
	"gorm.io/gorm"
)

type TableRepositoryInterface interface {
	ListTables(sourceId int, pagination dto.Query) ([]entity.Table, error)
	GetTableById(id int) (entity.Table, error)
	GetCount() int64
}

type TableRepository struct {
	DB *gorm.DB
}

func (repository *TableRepository) ListTables(sourceId int, pagination dto.Query) ([]entity.Table, error) {
	var tables []entity.Table
	offset := (pagination.Page - 1) * pagination.Limit
	if err := repository.DB.
		Offset(offset).
		Limit(pagination.Limit).
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

func (repository *TableRepository) GetCount() int64 {
	var count int64
	repository.DB.Model(&entity.Table{}).Count(&count)
	return count
}
