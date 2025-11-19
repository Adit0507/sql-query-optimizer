package parser

import (
	"fmt"
	"strconv"
)

type Parser struct {
	lexer     *Lexer
	curToken  Token
	peekToken Token
	errors    []string
}

func NewParser(input string) *Parser {
	l := NewLexer(input)
	p := &Parser{
		lexer:  l,
		errors: []string{},
	}

	// readin 2 tokens to initialize currtoken and peektoken
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) addError(msg string) {
	p.errors = append(p.errors, fmt.Sprintf("Line %d, Col %d: %s", p.curToken.Line, p.curToken.Column, msg))
}

func (p *Parser) curTokenIs(t TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.addError(fmt.Sprintf("expected %s, got %s", t, p.peekToken.Type))
	return false
}

func (p *Parser) Parse() Statement {
	if p.curTokenIs(SELECT) {
		return p.parseSelectStatement()
	}
	p.addError(fmt.Sprintf("unexpcted token %s", p.curToken.Type))

	return nil
}

func (p *Parser) parseSelectStatement() *SelectStatement {
	stmt := &SelectStatement{}

	if !(p.peekTokenIs(ASTERISK) || p.peekTokenIs(IDENT)) {
		p.addError("expected column name or '*'")
		return nil
	}

	// parsing select columsn
	p.nextToken()
	stmt.Columns = p.parseSelectColumns()

	if !p.expectPeek(FROM) {
		return nil
	}

	p.nextToken()
	stmt.From = p.parseTableRef()

	// parsing optional JOINs
	for p.peekTokenIs(JOIN) || p.peekTokenIs(INNER) || p.peekTokenIs(LEFT) || p.peekTokenIs(RIGHT) {
		p.nextToken()

		join := p.parseJoinClause()
		if join != nil {
			stmt.Joins = append(stmt.Joins, join)
		}
	}

	// parse optional WHERE clause
	if p.peekTokenIs(WHERE) {
		p.nextToken()
		p.nextToken() //movin to next column

		stmt.Where = p.parseExpression()
	}

	return stmt
}

func (p *Parser) parseExpression() Expression {
	return p.parseOrExpression()
}

func (p *Parser) parseOrExpression() Expression {
	left := p.parseAndExpression()

	for p.peekTokenIs(OR) {
		p.nextToken()

		op := p.curToken.Literal
		p.nextToken()
		right := p.parseAndExpression()
		left = &BinaryExpr{
			Left:     left,
			Operator: op,
			Right:    right,
		}
	}

	return left
}

func (p *Parser) parseAndExpression() Expression {
	left := p.parseComparisionExpression()

	for p.peekTokenIs(AND) {
		p.nextToken()
		op := p.curToken.Literal
		p.nextToken()
		right := p.parseComparisionExpression()
		left = &BinaryExpr{
			Left:     left,
			Operator: op,
			Right:    right,
		}
	}

	return left
}
func (p *Parser) parseComparisionExpression() Expression {
	left := p.parsePrimaryExpression()

	if p.peekTokenIs(EQ) || p.peekTokenIs(NEQ) || p.peekTokenIs(LT) || p.peekTokenIs(LTE) || p.peekTokenIs(GT) || p.peekTokenIs(GTE) {
		p.nextToken()
		op := p.curToken.Literal
		p.nextToken()
		right := p.parsePrimaryExpression()

		return &BinaryExpr{
			Left:     left,
			Operator: op,
			Right:    right,
		}
	}

	return left
}

func (p *Parser) parsePrimaryExpression() Expression {
	switch p.curToken.Type {
	case IDENT:
		return p.parseColumnRef()

	case INT:
		val, _ := strconv.Atoi(p.curToken.Literal)
		return &Literal{Type: IntLiteral, Value: val}

	case STRING:
		return &Literal{Type: StringLiteral, Value: p.curToken.Literal}
	case LPAREN:
		p.nextToken()
		expr := p.parseExpression()
		if !p.expectPeek(RPAREN) {
			return nil
		}

		return expr

	default:
		p.addError(fmt.Sprintf("unexpected token in expression: %s", p.curToken.Type))
		return nil
	}
}

func (p *Parser) parseJoinClause() *JoinClause {
	join := &JoinClause{}

	//join type
	if p.curTokenIs(INNER) {
		join.Type = InnerJoin
		if !p.expectPeek(JOIN) {
			return nil
		}
	} else if p.curTokenIs(LEFT) {
		join.Type = LeftJoin
		if !p.expectPeek(JOIN) {
			return nil
		}
	} else if p.curTokenIs(RIGHT) {
		join.Type = RightJoin
		if !p.expectPeek(JOIN) {
			return nil
		}
	} else {
		join.Type = InnerJoin
	}

	// parse table
	if !p.expectPeek(IDENT) {
		return nil
	}
	join.Table = p.parseTableRef()

	// parse on condition
	if !p.expectPeek(ON) {
		return nil
	}
	p.nextToken()
	join.Condition = p.parseExpression()

	return join
}

func (p *Parser) parseTableRef() *TableRef {
	table := &TableRef{
		Name: p.curToken.Literal,
	}

	if p.peekTokenIs(AS) {
		p.nextToken()
		if !p.expectPeek(IDENT) {
			return nil
		}

		table.Alias = p.curToken.Literal
	} else if p.peekTokenIs(IDENT) {
		p.nextToken()
		table.Alias = p.curToken.Literal
	}

	return table
}

func (p *Parser) parseColumnRef() Expression {
	col := &ColumnRef{
		Column: p.curToken.Literal,
	}

	// checking for table.colum syntax
	if p.peekTokenIs(DOT) {
		col.Table = col.Column
		p.nextToken()

		if !p.expectPeek(IDENT) && !p.expectPeek(ASTERISK) {
			return nil
		}
		if p.curTokenIs(ASTERISK) {
			return &StarExpr{Table: col.Table}
		}

		col.Column = p.curToken.Literal
	}

	return col
}

func (p *Parser) parseSelectColumns() []Expression {
	var cols []Expression

	if p.curTokenIs(ASTERISK) {
		cols = append(cols, &StarExpr{})
		return cols
	}

	// parsin first column
	cols = append(cols, p.parseColumnRef())

	for p.peekTokenIs(COMMA) {
		p.nextToken()
		p.nextToken()

		cols = append(cols, p.parseColumnRef())
	}

	return cols
}
