package entity

const PostgreSQLType = "postgresql"

type Source struct {
	ID       int
	Title    string
	Name     string
	Type     string //postgresql, mysql, oracle, innodb, etc
	Host     string
	Port     int
	Username string
	Password string
	Schema   string
	Alive    bool
	Tables   []Table
}
