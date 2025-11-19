package parser

import (
	"testing"
)

func TestLexer(t *testing.T) {
	input := `SELECT id, name FROM users WHERE id = 5`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{SELECT, "SELECT"},
		{IDENT, "id"},
		{COMMA, ","},
		{IDENT, "name"},
		{FROM, "FROM"},
		{IDENT, "users"},
		{WHERE, "WHERE"},
		{IDENT, "id"},
		{EQ, "="},
		{INT, "5"},
		{EOF, ""},
	}

	l := NewLexer(input)
	for i, ttt := range tests {
		tok := l.NextToken()

		if tok.Type != ttt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, ttt.expectedType, tok.Type)
		}

		if tok.Literal != ttt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, ttt.expectedLiteral, tok.Literal)
		}
	}
}

func TestParseSelectWithWhere(t *testing.T) {
	input := `SELECT id, name FROM users WHERE id = 5`

	p := NewParser(input)
	stmt := p.Parse()

	if len(p.Errors()) > 0 {
		t.Fatalf("parser has errors: %v", p.Errors())
	}

	selectStmt := stmt.(*SelectStatement)
	if selectStmt.Where == nil {
		t.Fatal("expected WHERE clause, got nil")
	}

	whereExpr, ok := selectStmt.Where.(*BinaryExpr)
	if !ok {
		t.Fatalf("WHERE is not BinaryExpr, got %T", selectStmt.Where)
	}

	if whereExpr.Operator != "=" {
		t.Fatalf("expected operator '=', got '%s'", whereExpr.Operator)
	}
}
