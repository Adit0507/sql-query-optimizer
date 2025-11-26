package plan

import "github.com/Adit0507/sql-query-optimizer/internal/catalog"

type LogicalPlan interface {
	Children() []LogicalPlan
	Schema() []catalog.Column
	String() string
}


