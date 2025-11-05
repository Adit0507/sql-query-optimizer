package parser

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF
	IDENT
	INT
	STRING

	// keywords
	SELECT
	FROM
	WHERE
	JOIN
	INNER
	LEFT
	RIGHT
	ON
	AND
	OR
	AS
	ORDER
	BY
	GROUP
	HAVING
	LIMIT
	OFFSET

	// operators
	EQ
	NEQ
	LT
	GT
	LTE
	GTE
	// delimiters
	COMMA
	SEMICOLON
	LPAREN
	RPAREN
	ASTERISK
	DOT
)

var keywords = map[string]TokenType{
	"SELECT": SELECT,
	"FROM":   FROM,
	"WHERE":  WHERE,
	"JOIN":   JOIN,
	"INNER":  INNER,
	"LEFT":   LEFT,
	"RIGHT":  RIGHT,
	"ON":     ON,
	"AND":    AND,
	"OR":     OR,
	"AS":     AS,
	"ORDER":  ORDER,
	"BY":     BY,
	"GROUP":  GROUP,
	"HAVING": HAVING,
	"LIMIT":  LIMIT,
	"OFFSET": OFFSET,
}

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

// checkin if identifier is keyword
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENT
}

// string representation of toekn type
func (t TokenType) String() string {
	switch t {
		case ILLEGAL:
		return "ILLEGAL"
	case EOF:
		return "EOF"
	case IDENT:
		return "IDENT"
	case INT:
		return "INT"
	case STRING:
		return "STRING"
	case SELECT:
		return "SELECT"
	case FROM:
		return "FROM"
	case WHERE:
		return "WHERE"
	case JOIN:
		return "JOIN"
	case INNER:
		return "INNER"
	case LEFT:
		return "LEFT"
	case RIGHT:
		return "RIGHT"
	case ON:
		return "ON"
	case AND:
		return "AND"
	case OR:
		return "OR"
	case AS:
		return "AS"
	case ORDER:
		return "ORDER"
	case BY:
		return "BY"
	case GROUP:
		return "GROUP"
	case HAVING:
		return "HAVING"
	case LIMIT:
		return "LIMIT"
	case OFFSET:
		return "OFFSET"
	case EQ:
		return "="
	case NEQ:
		return "!="
	case LT:
		return "<"
	case GT:
		return ">"
	case LTE:
		return "<="
	case GTE:
		return ">="
	case COMMA:
		return ","
	case SEMICOLON:
		return ";"
	case LPAREN:
		return "("
	case RPAREN:
		return ")"
	case ASTERISK:
		return "*"
	case DOT:
		return "."
	default:
		return "UNKNOWN"
	}
}