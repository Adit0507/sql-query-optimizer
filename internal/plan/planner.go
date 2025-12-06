package plan

import (
	"fmt"

	"github.com/Adit0507/sql-query-optimizer/internal/catalog"
	"github.com/Adit0507/sql-query-optimizer/internal/parser"
)

type Planner struct { // AST to logical plans
	catalog *catalog.Catalog
}

func NewPlanner(cat *catalog.Catalog) *Planner{
	return &Planner{catalog: cat}
}

func (p *Planner) CreateLogicalPlan(stmt parser.Statement) (LogicalPlan, error) {
	selectStmt, ok := stmt.(*parser.SelectStatement)
	if !ok{
		return  nil, fmt.Errorf("only SELECT statements supported")
	}

	return p.planSelect(selectStmt)
}

func (p *Planner) planSelect(stmt *parser.SelectStatement) (LogicalPlan, error){
	table, err := p.catalog.GetTable(stmt.From.Name)	//table scan
	if err != nil {
		return  nil, err
	}

	var plan LogicalPlan = &LogicalScan{
		TableName: stmt.From.Name,
		Table: table,
		Alias: stmt.From.Alias,
	}

	return plan, nil
}