package catalog

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
