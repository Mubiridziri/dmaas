package entity

type Field struct {
	ID       int
	Name     string
	Type     string
	Comment  string
	Nullable bool
	TableID  int
	Table    Table
}
