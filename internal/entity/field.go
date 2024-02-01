package entity

type Field struct {
	ID       int    `json:"-"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Comment  string `json:"comment"`
	Nullable bool   `json:"nullable"`
	TableID  int    `json:"-"`
}
