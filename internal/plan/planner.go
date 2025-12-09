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

func (p *Planner) convertExpr(expr parser.Expression) (Expr, error) {
	switch e := expr.(type) {
	case *parser.ColumnRef:
		return &ColumnExpr{
			Table:  e.Table,
			Column: e.Column,
		}, nil

	case *parser.Literal:
		var dataType catalog.DataType

		switch e.Type {
		case parser.IntLiteral:
			dataType = catalog.IntType
		case parser.StringLiteral:
			dataType = catalog.StringType

		}
		return &LiteralExpr{
			Value: e.Value,
			Type:  dataType,
		}, nil

	case *parser.BinaryExpr:
		left, err := p.convertExpr(e.Left)
		if err != nil {
			return nil, err
		}
		right, err := p.convertExpr(e.Right)
		if err != nil {
			return nil, err
		}
		return &BinaryExpr{
			Left:     left,
			Operator: e.Operator,
			Right:    right,
		}, nil

	default:
		return nil, fmt.Errorf("unsupported expression type: %T", expr)
	}
}
