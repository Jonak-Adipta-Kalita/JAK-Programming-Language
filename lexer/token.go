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
	TT_KEYWORD
	TT_INT
	TT_STRING
	TT_SEMICOLON
	TT_LPAREN
	TT_RPAREN
	TT_NEWLINE
)

var KEYWORDS = []string{
	"print",
}

type Token struct {
	Value string
	Type  TokenType
}
