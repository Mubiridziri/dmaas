package entity

type Table struct {
	ID       int
	Name     string
	Comment  string
	Fields   []Field
	SourceID int
}
