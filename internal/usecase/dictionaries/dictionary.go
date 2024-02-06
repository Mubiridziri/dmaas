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
	Title  string                              `json:"title"`
	Fields []CreateOrUpdateDictionaryFieldView `json:"fields"`
}

type CreateOrUpdateDictionaryFieldView struct {
	Title string `json:"title"`
	Name  string `json:"name"`
	Type  string `json:"type"`
}

type PaginatedDictionariesList struct {
	Total   int64            `json:"total"`
	Entries []DictionaryView `json:"entries"`
}

type Repository interface {
	CreateDictionary(dictionary *entity.Dictionary) error
	UpdateDictionary(dictionary *entity.Dictionary) error
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

	var fields []entity.DictionaryField

	for _, inputField := range input.Fields {
		field := entity.DictionaryField{
			Title: inputField.Title,
			Name:  inputField.Name,
			Type:  inputField.Type,
		}
		fields = append(fields, field)
	}
	dict.Fields = fields

	err := c.Repository.CreateDictionary(&dict)

	if err != nil {
		return DictionaryView{}, err
	}

	return fromDBDictionary(&dict), nil
}

//func (c Controller) UpdateDictionary(dictionary *entity.Dictionary) error {
//	return useCase.UpdateDictionary(dictionary)
//}

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
