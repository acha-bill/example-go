package pkg

import (
	"fmt"
	"sort"
)

var (
	ErrorNoTableName = fmt.Errorf("no table name provided")
	ErrorNoTable     = fmt.Errorf("no table found")
	ErrTableExists   = fmt.Errorf("table already exists")
)

// PrimaryKey is the type of the primary key of the model
type PrimaryKey int

// Model is the interface that all models must implement
type Model interface {
	GetID() PrimaryKey
	SetID(PrimaryKey)
}

// DB is the interface that all databases must implement
type DB interface {
	// Tables returns a slice of all the tables in the database
	Tables() []string
	// Table returns a table by its name
	Table(string) (Table, error)
	// AddTable adds a new table to the database
	AddTable(string) error
}

type db struct {
	tables map[string]Table
}

func (d *db) AddTable(s string) error {
	if s == "" {
		return ErrorNoTableName
	}
	if _, ok := d.tables[s]; ok {
		return ErrTableExists
	}
	d.tables[s] = &table{
		name: s,
		data: make(map[PrimaryKey]Model),
	}
	return nil
}

func (d *db) Tables() []string {
	var tables []string
	for name := range d.tables {
		tables = append(tables, name)
	}
	sort.Strings(tables)
	return tables
}

func (d *db) Table(s string) (Table, error) {
	if s == "" {
		return nil, ErrorNoTableName
	}
	if table, ok := d.tables[s]; ok {
		return table, nil
	}
	return nil, ErrorNoTable
}

// NewDB returns a new database
func NewDB() DB {
	return &db{
		tables: make(map[string]Table),
	}
}

var _ DB = &db{}
