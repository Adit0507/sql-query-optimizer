package parser

import (
	"strings"
	"unicode"
)

type Lexer struct {
	input        string
	position     int  //current position in input
	readPosition int  //reading position(after current chr)
	ch           byte //current char
	line         int
	column       int
}

func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}

	l.readChar()

	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
	l.column++

	if l.ch == '\n' {
		l.line++
		l.column = 0
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}

	return l.input[l.readPosition]
}

func (l *Lexer) NextToken() Token {
	var tok Token
	l.skipWhitespace()

	tok.Line = l.line
	tok.Column = l.column

	switch l.ch {
	case '=':
		tok = l.newToken(EQ, string(l.ch))
	case ',':
		tok = l.newToken(COMMA, string(l.ch))
	case ';':
		tok = l.newToken(SEMICOLON, string(l.ch))
	case '(':
		tok = l.newToken(LPAREN, string(l.ch))
	case ')':
		tok = l.newToken(RPAREN, string(l.ch))
	case '*':
		tok = l.newToken(ASTERISK, string(l.ch))
	case '.':
		tok = l.newToken(DOT, string(l.ch))

	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()

			tok = l.newToken(LTE, string(ch)+string(l.ch))
		} else if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(NEQ, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(LT, string(l.ch))
		}

	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(GTE, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(GT, string(l.ch))
		}

	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(NEQ, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(ILLEGAL, string(l.ch))
		}
	case '"', '\'':
		tok.Type = STRING
		tok.Literal = l.readString(l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdent(strings.ToUpper(tok.Literal))
			return tok
		} else if isDigit(l.ch) {
			tok.Type = INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = l.newToken(ILLEGAL, string(l.ch))
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) readString(quote byte) string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == quote || l.ch == 0 {
			break
		}
	}

	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) newToken(tokenType TokenType, literal string) Token {
	return Token{
		Type:    tokenType,
		Literal: literal,
		Line:    l.line,
		Column:  l.column,
	}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch))
}
func isDigit(ch byte) bool {
	return unicode.IsDigit(rune(ch))
}
