package lexer

import (
	"testing"

	"github.com/Jonak-Adipta-Kalita/JAK-Programming-Language/token"
)

func TestNextToken(t *testing.T) {
	input := `=+-!*/(){},;`

	tests := []struct {
		expectedType token.TokenType
		expectedLiteral string
	}{
		{token.TT_ASSIGN, "="},
		{token.TT_PLUS, "+"},
		{token.TT_MINUS, "-"},
		{token.TT_BANG, "!"},
		{token.TT_ASTER, "*"},
		{token.TT_SLASH, "/"},
		{token.TT_LPAREN, "("},
		{token.TT_RPAREN, ")"},
		{token.TT_LBRACE, "{"},
		{token.TT_RBRACE, "}"},
		{token.TT_COMMA, ","},
		{token.TT_SEMICOLON, ";"},
		{token.TT_EOF, ""},
	}
	
	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
			i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
			i, tt.expectedLiteral, tok.Literal)
		}
	}
}
