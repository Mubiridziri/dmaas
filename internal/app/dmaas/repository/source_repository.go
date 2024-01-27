package repository

import (
	"dmaas/internal/app/dmaas/entity"
	"gorm.io/gorm"
)

type SourceRepositoryInterface interface {
	CreateSource(source *entity.Source) error
	UpdateSource(source *entity.Source) error
	RemoveSource(source *entity.Source) error
	ListSources(page, limit int) ([]entity.Source, error)
	GetSourceById(id int) (entity.Source, error)
	GetCount() int64
}

type SourceRepository struct {
	DB *gorm.DB
}

func (repository *SourceRepository) CreateSource(source *entity.Source) error {
	return repository.DB.Create(source).Error
}

func (repository *SourceRepository) UpdateSource(source *entity.Source) error {
	return repository.DB.Save(source).Error
}

func (repository *SourceRepository) RemoveSource(source *entity.Source) error {
	return repository.DB.Delete(source).Error
}

func (repository *SourceRepository) ListSources(page, limit int) ([]entity.Source, error) {
	var sources []entity.Source
	offset := (page - 1) * limit
	if err := repository.DB.Offset(offset).Limit(limit).Find(&sources).Error; err != nil {
		return []entity.Source{}, err
	}

	return sources, nil
}

func (repository *SourceRepository) GetSourceById(id int) (entity.Source, error) {
	var source entity.Source
	if err := repository.DB.
		Preload("Tables").
		Preload("Tables.Fields").
		Where(entity.Source{ID: id}).
		First(&source).Error; err != nil {
		return entity.Source{}, err
	}

	return source, nil
}

func (repository *SourceRepository) GetCount() int64 {
	var count int64
	repository.DB.Model(&entity.Source{}).Count(&count)
	return count
}
