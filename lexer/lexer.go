package lexer

import (
	"fmt"

	"github.com/Jonak-Adipta-Kalita/JAK-Programming-Language/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
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
	l.readPosition += 1
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	fmt.Println(string(l.ch))

	switch l.ch {
		case '=':
			tok = newToken(token.TT_ASSIGN, l.ch)
		case ';':
			tok = newToken(token.TT_SEMICOLON, l.ch)
		case '(':
			tok = newToken(token.TT_LPAREN, l.ch)
		case ')':
			tok = newToken(token.TT_RPAREN, l.ch)
		case ',':
			tok = newToken(token.TT_COMMA, l.ch)
		case '+':
			tok = newToken(token.TT_PLUS, l.ch)
		case '-':
			tok = newToken(token.TT_MINUS, l.ch)
		case '!':
			tok = newToken(token.TT_BANG, l.ch)
		case '*':
			tok = newToken(token.TT_ASTER, l.ch)
		case '/':
			tok = newToken(token.TT_SLASH, l.ch)
		case '{':
			tok = newToken(token.TT_LBRACE, l.ch)
		case '}':
			tok = newToken(token.TT_RBRACE, l.ch)
		case 0:
			tok.Literal = ""
			tok.Type = token.TT_EOF
	}

	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}