package dictionaries

import (
	"dmaas/internal/entity"
	"time"
)

type DictionaryView struct {
	ID        int                   `json:"id"`
	Title     string                `json:"title"`
	Fields    []DictionaryFieldView `json:"fields"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
}

type DictionaryFieldView struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Name  string `json:"name"`
	Type  string `json:"type"`
}

type CreateOrUpdateDictionaryView struct {
	Title  string                              `json:"title"  binding:"required"`
	Fields []CreateOrUpdateDictionaryFieldView `json:"fields"  binding:"required"`
}

type CreateOrUpdateDictionaryFieldView struct {
	ID    int    `json:"id"`
	Title string `json:"title"  binding:"required"`
	Name  string `json:"name"  binding:"required"`
	Type  string `json:"type"  binding:"required"`
}

type PaginatedDictionariesList struct {
	Total   int64            `json:"total"`
	Entries []DictionaryView `json:"entries"`
}

type Repository interface {
	CreateDictionary(dictionary *entity.Dictionary) error
	UpdateDictionary(dictionary *entity.Dictionary, removeFields []entity.DictionaryField) error
	RemoveDictionary(dictionary *entity.Dictionary) error
	ListDictionaries(page, limit int) ([]entity.Dictionary, error)
	GetDictionaryById(id int) (entity.Dictionary, error)
	GetDictionariesCount() int64
}

type Controller struct {
	Repository
}

func NewController(repo Repository) *Controller {
	return &Controller{Repository: repo}
}

func (c Controller) CreateDictionary(input CreateOrUpdateDictionaryView) (DictionaryView, error) {
	dict := entity.Dictionary{
		Title: input.Title,
	}

	for _, inputField := range input.Fields {
		field := entity.DictionaryField{
			Title: inputField.Title,
			Name:  inputField.Name,
			Type:  inputField.Type,
		}
		dict.Fields = append(dict.Fields, field)
	}

	err := c.Repository.CreateDictionary(&dict)

	if err != nil {
		return DictionaryView{}, err
	}

	return fromDBDictionary(&dict), nil
}

func (c Controller) UpdateDictionary(id int, input CreateOrUpdateDictionaryView) (DictionaryView, error) {
	dict, err := c.Repository.GetDictionaryById(id)

	if err != nil {
		return DictionaryView{}, err
	}

	dict.Title = input.Title

	//adding & updating fields
	for _, inputField := range input.Fields {
		exists := false
		for _, field := range dict.Fields {
			if field.ID == inputField.ID {
				field.Name = inputField.Name
				field.Title = inputField.Title
				field.Type = inputField.Type
				exists = true
				break
			}
		}

		if !exists {
			dict.Fields = append(dict.Fields, entity.DictionaryField{
				Title: inputField.Title,
				Name:  inputField.Name,
				Type:  inputField.Type,
			})
		}
	}

	var removeFields []entity.DictionaryField

	for index, field := range dict.Fields {
		exists := false
		for _, inputField := range input.Fields {
			if field.ID == inputField.ID {
				exists = true
				break
			}
		}

		if !exists {
			dict.Fields = append(dict.Fields[:index], dict.Fields[index+1:]...)
			removeFields = append(removeFields, field)
		}
	}

	err = c.Repository.UpdateDictionary(&dict, removeFields)

	if err != nil {
		return DictionaryView{}, err
	}

	return fromDBDictionary(&dict), nil

}

func (c Controller) RemoveDictionary(id int) (DictionaryView, error) {

	dict, err := c.Repository.GetDictionaryById(id)

	if err != nil {
		return DictionaryView{}, err
	}

	err = c.Repository.RemoveDictionary(&dict)

	return fromDBDictionary(&dict), err
}

func (c Controller) ListDictionaries(page, limit int) (PaginatedDictionariesList, error) {
	dicts, err := c.Repository.ListDictionaries(page, limit)

	if err != nil {
		return PaginatedDictionariesList{}, err
	}

	var entries []DictionaryView

	for _, dict := range dicts {
		entries = append(entries, fromDBDictionary(&dict))
	}

	return PaginatedDictionariesList{
		Total:   c.Repository.GetDictionariesCount(),
		Entries: entries,
	}, nil
}

func (c Controller) GetDictionaryById(id int) (DictionaryView, error) {
	dict, err := c.Repository.GetDictionaryById(id)

	return fromDBDictionary(&dict), err
}

func fromDBDictionary(dictionary *entity.Dictionary) DictionaryView {
	dictionaryView := DictionaryView{
		ID:        dictionary.ID,
		Title:     dictionary.Title,
		CreatedAt: dictionary.CreatedAt,
		UpdatedAt: dictionary.UpdatedAt,
	}

	var fieldsViews []DictionaryFieldView

	for _, field := range dictionary.Fields {
		fieldView := DictionaryFieldView{
			ID:    field.ID,
			Title: field.Title,
			Name:  field.Name,
			Type:  field.Type,
		}
		fieldsViews = append(fieldsViews, fieldView)
	}
	dictionaryView.Fields = fieldsViews

	return dictionaryView
}
