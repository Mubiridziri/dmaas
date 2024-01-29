package entity

import "time"

type Dictionary struct {
	ID        int               `json:"id"`
	Title     string            `json:"title"`
	Fields    []DictionaryField `json:"fields"`
	Rows      []DictionaryRow   `json:"-"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}
