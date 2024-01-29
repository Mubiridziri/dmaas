package repository

import (
	"dmaas/internal/app/dmaas/dto"
	"dmaas/internal/app/dmaas/entity"
	"gorm.io/gorm"
)

type DictionaryRepositoryInterface interface {
	CreateDictionary(dictionary *entity.Dictionary) error
	UpdateDictionary(dictionary *entity.Dictionary) error
	RemoveDictionary(dictionary *entity.Dictionary) error
	ListDictionaries(pagination dto.Query) ([]entity.Dictionary, error)
	GetDictionaryById(id int) (entity.Dictionary, error)
	GetCount() int64
}

type DictionaryRepository struct {
	DB *gorm.DB
}

func (repository *DictionaryRepository) CreateDictionary(dictionary *entity.Dictionary) error {
	//Сохраняем справочник в транзакции с полями, если возникает ошибка откатываем всё, включая справочник
	err := repository.DB.Transaction(func(tx *gorm.DB) error {
		err := repository.DB.Create(dictionary).Error
		if err != nil {
			return err
		}

		for _, field := range dictionary.Fields {
			field.DictionaryID = dictionary.ID
			err = repository.DB.Create(&field).Error
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func (repository *DictionaryRepository) UpdateDictionary(dictionary *entity.Dictionary) error {
	//Сохраняем справочник в транзакции с полями, если возникает ошибка откатываем всё, включая справочник
	err := repository.DB.Transaction(func(tx *gorm.DB) error {
		err := repository.DB.Save(dictionary).Error
		if err != nil {
			return err
		}

		for _, field := range dictionary.Fields {
			field.DictionaryID = dictionary.ID
			if field.ID == 0 {
				err = repository.DB.Create(&field).Error
			} else {
				err = repository.DB.Save(&field).Error
			}
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err

}

func (repository *DictionaryRepository) RemoveDictionary(dictionary *entity.Dictionary) error {
	return repository.DB.Delete(dictionary).Error
}

func (repository *DictionaryRepository) ListDictionaries(pagination dto.Query) ([]entity.Dictionary, error) {
	var dictionaries []entity.Dictionary
	offset := (pagination.Page - 1) * pagination.Limit
	if err := repository.DB.Offset(offset).Limit(pagination.Limit).Find(&dictionaries).Error; err != nil {
		return []entity.Dictionary{}, err
	}

	return dictionaries, nil
}

func (repository *DictionaryRepository) GetDictionaryById(id int) (entity.Dictionary, error) {
	var dictionary entity.Dictionary
	if err := repository.DB.
		Preload("Fields").
		Where(entity.Dictionary{ID: id}).
		First(&dictionary).Error; err != nil {
		return entity.Dictionary{}, err
	}

	return dictionary, nil
}

func (repository *DictionaryRepository) GetCount() int64 {
	var count int64
	repository.DB.Model(&entity.Dictionary{}).Count(&count)
	return count
}
