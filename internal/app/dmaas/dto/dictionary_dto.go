package dto

import "dmaas/internal/app/dmaas/entity"

type DictionaryRequest struct {
	Title  string                   `json:"title"`
	Fields []DictionaryFieldRequest `json:"fields"`
}

type DictionaryFieldRequest struct {
	Title string `json:"title" binding:"required"`
	Name  string `json:"name" binding:"required"`
	Type  string `json:"type" binding:"required"`
}

type DictionaryUpdateRequest struct {
	Title  string                         `json:"title" binding:"required"`
	Fields []DictionaryFieldUpdateRequest `json:"fields" binding:"required"`
}

type DictionaryFieldUpdateRequest struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Name  string `json:"name"`
	Type  string `json:"type"`
}

func (dto *DictionaryRequest) ToDictionary() entity.Dictionary {
	dictionary := entity.Dictionary{
		Title: dto.Title,
	}

	var fields []entity.DictionaryField

	for _, fieldRequest := range dto.Fields {
		field := entity.DictionaryField{
			Title: fieldRequest.Title,
			Name:  fieldRequest.Name,
			Type:  fieldRequest.Type,
		}
		fields = append(fields, field)
	}

	dictionary.Fields = fields
	return dictionary
}

func (dto *DictionaryUpdateRequest) ToDictionary(dictionary *entity.Dictionary) {
	dictionary.Title = dto.Title
	fields := dictionary.Fields

	for _, updatedField := range dto.Fields {
		found := false
		for _, oldField := range fields {
			if updatedField.ID == oldField.ID {
				oldField.Name = updatedField.Name
				oldField.Title = updatedField.Title
				oldField.Type = updatedField.Type
				found = true
				break
			}
		}
		if !found {
			field := entity.DictionaryField{
				DictionaryID: dictionary.ID,
				Title:        updatedField.Title,
				Name:         updatedField.Name,
				Type:         updatedField.Type,
			}
			fields = append(fields, field)
		}
	}

	dictionary.Fields = fields
}
