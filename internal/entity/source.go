package entity

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

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

type Table struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Comment  string  `json:"comment"`
	Fields   []Field `json:"fields"`
	SourceID int     `json:"-"`
}

type Field struct {
	ID       int    `json:"-"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Comment  string `json:"comment"`
	Nullable bool   `json:"nullable"`
	TableID  int    `json:"-"`
}

type sourceRepository struct {
	db *gorm.DB
}

type tableRepository struct {
	db *gorm.DB
}

type tableDataRepository struct {
	db *gorm.DB
}

func (r sourceRepository) CreateSource(source *Source) error {
	return r.db.Create(source).Error
}

func (r sourceRepository) UpdateSource(source *Source) error {
	return r.db.Save(source).Error
}

func (r sourceRepository) RemoveSource(source *Source) error {
	return r.db.Delete(source).Error
}

func (r sourceRepository) ListSources(page, limit int) ([]Source, error) {
	var sources []Source
	offset := (page - 1) * limit
	if err := r.db.Offset(offset).Limit(limit).Find(&sources).Error; err != nil {
		return []Source{}, err
	}

	return sources, nil
}

func (r sourceRepository) GetSourceById(id int) (Source, error) {
	var source Source
	if err := r.db.
		Preload("Tables").
		Preload("Tables.Fields").
		Where(Source{ID: id}).
		First(&source).Error; err != nil {
		return Source{}, err
	}

	return source, nil
}

func (r sourceRepository) GetSourcesCount() int64 {
	var count int64
	r.db.Model(&Source{}).Count(&count)
	return count
}

func (r tableRepository) ListTables(source Source, page, limit int) ([]Table, error) {
	var tables []Table
	offset := (page - 1) * limit
	if err := r.db.
		Offset(offset).
		Limit(limit).
		Where(Table{SourceID: source.ID}).
		Find(&tables).Error; err != nil {
		return []Table{}, err
	}

	return tables, nil
}

func (r tableRepository) GetTableById(id int) (Table, error) {
	var table Table
	if err := r.db.
		Preload("Fields").
		Where(Table{ID: id}).
		First(&table).Error; err != nil {
		return Table{}, err
	}

	return table, nil
}

func (r tableRepository) GetTablesCount() int64 {
	var count int64
	r.db.Model(&Table{}).Count(&count)
	return count
}

func (r tableDataRepository) ListTableData(localSchemaName string, table Table, page, limit int) ([]map[string]interface{}, error) {
	offset := (page - 1) * limit
	var data []map[string]interface{}
	sql := fmt.Sprintf("SELECT * FROM %v.%v OFFSET %v LIMIT %v", localSchemaName, table.Name, offset, limit)
	err := r.db.Raw(sql).Scan(&data).Error

	return data, err
}

func (r tableDataRepository) GetTableDataCount(localSchemaName string, table Table) int64 {
	var data map[string]int64
	sql := fmt.Sprintf("SELECT COUNT(*) as count FROM %v.%v", localSchemaName, table.Name)
	err := r.db.Raw(sql).Scan(&data).Error

	if err != nil {
		return 0
	}

	return data["count"]
}
