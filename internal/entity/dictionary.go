package entity

import (
	"gorm.io/gorm"
	"time"
)

const (
	TextFieldType     = "text"
	NumberFieldType   = "number"
	BooleanFieldType  = "boolean"
	DateFieldType     = "date"
	DateTimeFieldType = "datetime"
)

type Dictionary struct {
	ID        int               `json:"id"`
	Title     string            `json:"title"`
	Fields    []DictionaryField `json:"fields"`
	Rows      []DictionaryRow   `json:"-"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type DictionaryField struct {
	ID           int               `json:"id"`
	DictionaryID int               `json:"-"`
	Title        string            `json:"title"`
	Name         string            `json:"name"`
	Type         string            `json:"type"`
	Values       []DictionaryValue `json:"-"`
}

type DictionaryRow struct {
	ID           int               `json:"id"`
	DictionaryID int               `json:"-"`
	Values       []DictionaryValue `json:"values"`
}

type DictionaryValue struct {
	ID                int    `json:"id"`
	DictionaryFieldID int    `json:"-"`
	DictionaryRowID   int    `json:"-"`
	Value             string `json:"value"`
}

type dictionaryRepository struct {
	db *gorm.DB
}

type dictionaryDataRepository struct {
	db *gorm.DB
}

func (r dictionaryRepository) CreateDictionary(dictionary *Dictionary) error {
	return r.db.Create(dictionary).Error
}

func (r dictionaryRepository) UpdateDictionary(dictionary *Dictionary, removeFields []DictionaryField) error {
	//Сохраняем справочник в транзакции с полями, если возникает ошибка откатываем всё, включая справочник
	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := r.db.Save(dictionary).Error
		if err != nil {
			return err
		}

		for _, field := range removeFields {
			err = r.db.Delete(&field).Error
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err

}

func (r dictionaryRepository) RemoveDictionary(dictionary *Dictionary) error {
	return r.db.Delete(dictionary).Error
}

func (r dictionaryRepository) ListDictionaries(page, limit int) ([]Dictionary, error) {
	var dictionaries []Dictionary
	offset := (page - 1) * limit
	if err := r.db.Offset(offset).Limit(limit).Find(&dictionaries).Error; err != nil {
		return []Dictionary{}, err
	}

	return dictionaries, nil
}

func (r dictionaryRepository) GetDictionaryById(id int) (Dictionary, error) {
	var dictionary Dictionary
	if err := r.db.
		Preload("Fields").
		Where(Dictionary{ID: id}).
		First(&dictionary).Error; err != nil {
		return Dictionary{}, err
	}

	return dictionary, nil
}

func (r dictionaryRepository) GetDictionariesCount() int64 {
	var count int64
	r.db.Model(&Dictionary{}).Count(&count)
	return count
}

func (r dictionaryDataRepository) ListDictionaryData(dictionary Dictionary, page, limit int) ([]DictionaryRow, error) {
	var rows []DictionaryRow
	offset := (page - 1) * limit
	if err := r.db.
		Where(DictionaryRow{DictionaryID: dictionary.ID}).
		Preload("Values").
		Offset(offset).
		Limit(limit).
		Find(&rows).Error; err != nil {
		return []DictionaryRow{}, err
	}

	return rows, nil
}

func (r dictionaryDataRepository) GetDictionariesDataCount(dictionary Dictionary) int64 {
	var count int64
	r.db.Model(&DictionaryRow{}).Where(DictionaryRow{DictionaryID: dictionary.ID}).Count(&count)
	return count
}
