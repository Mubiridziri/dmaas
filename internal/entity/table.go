package entity

type Table struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Comment  string  `json:"comment"`
	Fields   []Field `json:"fields"`
	SourceID int     `json:"-"`
}
