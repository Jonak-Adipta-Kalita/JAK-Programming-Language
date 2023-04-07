package lexer

type TokenType uint

const (
	TT_EOF TokenType = iota
	TT_ILLEGAL
	TT_WHITE_SPACE
	TT_FLOAT
	TT_INT
	TT_STRING
	TT_BOOL
	TT_PLUS
	TT_MINUS
	TT_MUL
	TT_DIV
	TT_POW
	TT_ASSIGN
	TT_LPAREN
	TT_RPAREN
	TT_LSQUARE
	TT_RSQUARE
)

type Token struct {
	Value string
	Type  TokenType
}