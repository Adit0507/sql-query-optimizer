package parser

type Parser struct {
	lexer     *Lexer
	curToken  Token
	peekToken Token
	errors    []string
}

func NewParser(input string) *Parser {
	l := NewLexer(input)
	p := &Parser{
		lexer: l,
		errors: []string{},
	} 

	// readin 2 tokens to initialize currtoken and peektoken
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken () {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) Errors() []string {
	return p.errors
}


