package parser

// basic interface of ast nodes
type Node interface {
	String() string
}

type Statement interface { //sql statement
	Node
	statementNode()
}

type Expression interface { //sql expression
	Node
	expressionNode()
}

// SELECT statement
type SelectStatement struct {
	Columns []Expression
	From    *TableRef
	Where   Expression
	Joins   []*JoinClause
	OrderBy []*OrderByExpr
	GroupBy []Expression
	Limit   *int
	Offset  *int
}

func (s *SelectStatement) statementNode() {}
func (s *SelectStatement) String() string {
	return "SELECT"
}

func (t *TableRef) expressionNode() {}
func (t *TableRef) String() string {
	return t.Name
}

type JoinType int
type TableRef struct { //repersents table ref
	Name  string
	Alias string
}

type JoinClause struct {
	Type      JoinType
	Table     *TableRef
	Condition Expression
}

const (
	InnerJoin JoinType = iota
	LeftJoin
	RightJoin
)

func (j *JoinClause) String() string { return "JOIN" }

type OrderByExpr struct { //ORDER BY expression
	Expr Expression
	Desc bool
}

type ColumnRef struct { // col reference
	Table  string
	Column string
}

func (c *ColumnRef) expressionNode() {}

func (c *ColumnRef) String() string {
	if c.Table != "" {
		return c.Table + "." + c.Column
	}

	return c.Column
}
