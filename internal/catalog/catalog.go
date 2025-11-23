package catalog

import (
	"encoding/json"
	"fmt"
	"os"
)

type DataType int //simple col types

const (
	IntType DataType = iota
	StringType
	BoolType
)

func (d DataType) String() string {
	switch d {
	case IntType:
		return "INT"
	case BoolType:
		return "BOOL"
	case StringType:
		return "STRING"

	default:
		return "UNKNOWN"
	}
}

type Column struct { //table column
	Name string   `json:"name"`
	Type DataType `json:"type"`
}

type Index struct { //table index
	Name    string   `json:"name"`
	Columns []string `json:"columns"`
}

type Statistics struct { //table statistics for cost estimation
	RowCount      int            `json:"row_count"`
	DistinctCount map[string]int `json:"distinct_count"`
	NullCount     map[string]int `json:"null_count"`
}

type TableInfo struct { //metadata about table
	Name       string      `json:"name"`
	Columns    []Column    `json:"columns"`
	Indexes    []Index     `json:"indexes"`
	Statistics *Statistics `json:"statistics"`
	DataFile   string      `json:"data_file"` //pat to json datafile
}

type Catalog struct {
	tables map[string]*TableInfo
}

func NewCatalog() *Catalog {
	return &Catalog{
		tables: make(map[string]*TableInfo),
	}
}

// loading metdata from json fle
func (c *Catalog) LoadFromFile(filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read catalog file: %w", err)
	}

	var tables []TableInfo
	if err := json.Unmarshal(data, &tables); err != nil {
		return fmt.Errorf("failed to parse catalog %w", err)
	}

	for i := range tables {
		c.tables[tables[i].Name] = &tables[i]
	}

	return nil
}

// adds a table to catalog
func (c *Catalog) RegisterTable(table *TableInfo) {
	c.tables[table.Name] = table
}
