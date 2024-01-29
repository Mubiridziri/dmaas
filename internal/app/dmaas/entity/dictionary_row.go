package entity

type DictionaryRow struct {
	ID           int               `json:"id"`
	DictionaryID int               `json:"-"`
	Values       []DictionaryValue `json:"values"`
}
