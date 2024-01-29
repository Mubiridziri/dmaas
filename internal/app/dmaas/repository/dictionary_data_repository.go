package repository

import (
	"dmaas/internal/app/dmaas/dto"
	"dmaas/internal/app/dmaas/entity"
	"gorm.io/gorm"
)

type DictionaryDataRepositoryInterface interface {
	ListDictionaryData(dictionary entity.Dictionary, pagination dto.Query) ([]entity.DictionaryRow, error)
	GetCount(dictionary entity.Dictionary) int64
}

type DictionaryDataRepository struct {
	DB *gorm.DB
}

func (repository *DictionaryDataRepository) ListDictionaryData(
	dictionary entity.Dictionary,
	pagination dto.Query,
) ([]entity.DictionaryRow, error) {
	var rows []entity.DictionaryRow
	offset := (pagination.Page - 1) * pagination.Limit
	if err := repository.DB.
		Where(entity.DictionaryRow{DictionaryID: dictionary.ID}).
		Preload("Values").
		Offset(offset).
		Limit(pagination.Limit).
		Find(&rows).Error; err != nil {
		return []entity.DictionaryRow{}, err
	}

	return rows, nil
}

func (repository *DictionaryDataRepository) GetCount(dictionary entity.Dictionary) int64 {
	var count int64
	repository.DB.Model(&entity.DictionaryRow{}).Where(entity.DictionaryRow{DictionaryID: dictionary.ID}).Count(&count)
	return count
}
