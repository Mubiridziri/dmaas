package dictionarydata

import (
	"dmaas/internal/entity"
)

type DictionaryRowView struct {
	Values []DictionaryValueView
}

type DictionaryValueView struct {
	Field string
	Value string
}

type PaginatedDictionaryDataList struct {
	Total   int64               `json:"total"`
	Entries []DictionaryRowView `json:"entries"`
}

type Repository interface {
	ListDictionaryData(dictionary entity.Dictionary, page, limit int) ([]entity.DictionaryRow, error)
	GetDictionariesDataCount(dictionary entity.Dictionary) int64
}

type Controller struct {
	Repository
}

func NewController(repo Repository) *Controller {
	return &Controller{Repository: repo}
}

func (c Controller) ListDictionaryData(dictionary entity.Dictionary, page, limit int) (PaginatedDictionaryDataList, error) {

	rows, err := c.Repository.ListDictionaryData(dictionary, page, limit)

	if err != nil {
		return PaginatedDictionaryDataList{}, err
	}

	var entries []DictionaryRowView

	for _, row := range rows {
		entries = append(entries, fromDBDictionaryValue(&row))
	}

	return PaginatedDictionaryDataList{
		Total:   c.Repository.GetDictionariesDataCount(dictionary),
		Entries: entries,
	}, nil
}

func fromDBDictionaryValue(row *entity.DictionaryRow) DictionaryRowView {
	var values []DictionaryValueView

	for _, value := range row.Values {
		valueView := DictionaryValueView{
			Field: "Field", //TODO value.Field.Title
			Value: value.Value,
		}
		values = append(values, valueView)
	}

	return DictionaryRowView{Values: values}
}
