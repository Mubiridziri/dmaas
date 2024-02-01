package entity

const (
	TextFieldType     = "text"
	NumberFieldType   = "number"
	BooleanFieldType  = "boolean"
	DateFieldType     = "date"
	DateTimeFieldType = "datetime"
)

type DictionaryField struct {
	ID           int               `json:"id"`
	DictionaryID int               `json:"-"`
	Title        string            `json:"title"`
	Name         string            `json:"name"`
	Type         string            `json:"type"`
	Values       []DictionaryValue `json:"-"`
}
