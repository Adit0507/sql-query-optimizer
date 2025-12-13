package plan

import (
	"fmt"

	"github.com/Adit0507/sql-query-optimizer/internal/catalog"
	"github.com/Adit0507/sql-query-optimizer/internal/parser"
)

type Planner struct { // AST to logical plans
	catalog *catalog.Catalog
}

func NewPlanner(cat *catalog.Catalog) *Planner {
	return &Planner{catalog: cat}
}

func (p *Planner) CreateLogicalPlan(stmt parser.Statement) (LogicalPlan, error) {
	selectStmt, ok := stmt.(*parser.SelectStatement)
	if !ok {
		return nil, fmt.Errorf("only SELECT statements supported")
	}

	return p.planSelect(selectStmt)
}

func (p *Planner) planSelect(stmt *parser.SelectStatement) (LogicalPlan, error) {
	table, err := p.catalog.GetTable(stmt.From.Name) //table scan
	if err != nil {
		return nil, err
	}

	var plan LogicalPlan = &LogicalScan{
		TableName: stmt.From.Name,
		Table:     table,
		Alias:     stmt.From.Alias,
	}

	// joins
	for _, join := range stmt.Joins {
		rightTable, err := p.catalog.GetTable(join.Table.Name)
		if err != nil {
			return nil, err
		}

		rightScan := &LogicalScan{
			TableName: join.Table.Name,
			Table:     rightTable,
			Alias:     join.Table.Alias,
		}

		condition, err := p.convertExpr(join.Condition)
		if err != nil {
			return nil, err
		}

		plan = &LogicalJoin{
			Left:      plan,
			Right:     rightScan,
			JoinType:  convertJoinType(join.Type),
			Condition: condition,
		}
	}

	// add WEHERE filter
	if stmt.Where != nil {
		predicate, err := p.convertExpr(stmt.Where)
		if err != nil {
			return nil, err
		}

		plan = &LogicalFilter{
			Input:     plan,
			Predicate: predicate,
		}
	}

	// addin prjections
	projections, columnNames, err := p.convertProjections(stmt.Columns, plan)
	if err != nil {
		return nil, err
	}

	plan = &LogicalProject{
		Input:       plan,
		Projections: projections,
		ColumnNames: columnNames,
	}

	return plan, nil
}

func (p *Planner) convertProjections(cols []parser.Expression, input LogicalPlan) ([]Expr, []string, error) {
	var projections []Expr
	var columnNames []string

	for _, col := range cols {
		switch c := col.(type) {
		case *parser.StarExpr:
			schema := input.Schema()
			for _, schemaCol := range schema {
				projections = append(projections, &ColumnExpr{
					Column: schemaCol.Name,
				})
				columnNames = append(columnNames, schemaCol.Name)
			}

		case *parser.ColumnRef:
			projections = append(projections, &ColumnExpr{
				Table:  c.Table,
				Column: c.Column,
			})
			columnNames = append(columnNames, c.Column)

		default:
			return nil, nil, fmt.Errorf("unsupported projection type: %T", col)
		}
	}

	return projections, columnNames, nil
}

func convertJoinType(jt parser.JoinType) JoinType {
	switch jt {
	case parser.InnerJoin:
		return InnerJoin
	case parser.LeftJoin:
		return LeftJoin
	case parser.RightJoin:
		return RightJoin
	default:
		return InnerJoin
	}
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

func PrintPlan(plan LogicalPlan, indent int) {
	prefix := ""

	for i := 0; i < indent; i++ {
		prefix += "  "
	}

	fmt.Printf("%s%s\n", prefix, plan.String())
	for _, child := range plan.Children() {
		PrintPlan(child, indent+1)
	}
}