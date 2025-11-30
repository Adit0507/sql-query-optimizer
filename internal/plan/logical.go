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
	return  l.Table.Columns
}

func (l *LogicalScan) String() string {
	if l.Alias != ""{
		return fmt.Sprintf("Scan(%s AS %s)", l.TableName, l.Alias)
	}	

	return fmt.Sprintf("Scan(%s)", l.TableName)
}
