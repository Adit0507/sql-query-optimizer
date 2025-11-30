package plan

import (
	"fmt"

	"github.com/Adit0507/sql-query-optimizer/internal/catalog"
)

type LogicalPlan interface {
	Children() []LogicalPlan
	Schema() []catalog.Column
	String() string
}

// representin table scan
type LogicalScan struct {
	TableName string
	Table     *catalog.TableInfo
	Alias     string
}

func (l *LogicalScan) Children() []LogicalPlan {
	return nil
}
func (l *LogicalScan) Schema() []catalog.Column {
	return l.Table.Columns
}

func (l *LogicalScan) String() string {
	if l.Alias != "" {
		return fmt.Sprintf("Scan(%s AS %s)", l.TableName, l.Alias)
	}

	return fmt.Sprintf("Scan(%s)", l.TableName)
}

type LogicalFilter struct {
	Input     LogicalPlan
	Predicate Expr
}

func (l *LogicalFilter) Children() []LogicalPlan {
	return []LogicalPlan{l.Input}
}
func (l *LogicalFilter) Schema() []catalog.Column {
	return l.Input.Schema()
}
func (l *LogicalFilter) String() string {
	return fmt.Sprintf("Filter(%s)", l.Predicate.String())
}

// expression in logical plan
type Expr interface {
	String() string
}

// SELECT column list
type LogicalProject struct {
	Input       LogicalPlan
	Projections []Expr
	ColumnNames []string
}

func (l *LogicalProject) Children() []LogicalPlan {
	return []LogicalPlan{l.Input}
}

func (l *LogicalProject) Schema() []catalog.Column {
	cols := make([]catalog.Column, len(l.ColumnNames))

	for i, name := range l.ColumnNames { //columnnames based on projections
		cols[i] = catalog.Column{Name: name, Type: catalog.StringType}
	}

	return cols
}
func (l *LogicalProject) String() string {
	return fmt.Sprintf("Project(%v)", l.ColumnNames)
}
