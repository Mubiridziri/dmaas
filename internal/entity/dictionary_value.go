package entity

type DictionaryValue struct {
	ID                int    `json:"id"`
	DictionaryFieldID int    `json:"-"`
	DictionaryRowID   int    `json:"-"`
	Value             string `json:"value"`
}
