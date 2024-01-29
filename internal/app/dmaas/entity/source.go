package entity

import "time"

const PostgreSQLType = "postgresql"

type Source struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Name      string    `json:"name"`
	Type      string    `json:"type"` //postgresql, mysql, oracle, innodb, etc
	Host      string    `json:"host"`
	Port      int       `json:"port"`
	Username  string    `json:"-"`
	Password  string    `json:"-"`
	Schema    string    `json:"schema"`
	Alive     bool      `json:"alive"`
	Tables    []Table   `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
