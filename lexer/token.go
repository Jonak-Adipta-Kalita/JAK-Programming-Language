package lexer

type TokenType uint

const (
	TT_EOF TokenType = iota
	TT_ILLEGAL
	TT_WHITE_SPACE
	TT_PLUS
	TT_MINUS
	TT_MUL
	TT_DIV
	TT_POW
	TT_ASSIGN
	TT_LPAREN
	TT_RPAREN
)

type Token struct {
	Value string
	Type  TokenType
}
