package lexer

import (
	"testing"

	"github.com/Jonak-Adipta-Kalita/JAK-Programming-Language/token"
)

func TestNextToken(t *testing.T) {
	input := `print(((2 + 4 - 6 * 9) / 2) != 2);`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.PRINT, "print"},
		{token.LPAREN, "("},
		{token.LPAREN, "("},
		{token.LPAREN, "("},
		{token.INT, "2"},
		{token.PLUS, "+"},
		{token.INT, "4"},
		{token.MINUS, "-"},
		{token.INT, "6"},
		{token.ASTERISK, "*"},
		{token.INT, "9"},
		{token.RPAREN, ")"},
		{token.SLASH, "/"},
		{token.INT, "2"},
		{token.RPAREN, ")"},
		{token.NOT_EQ, "!="},
		{token.INT, "2"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
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
