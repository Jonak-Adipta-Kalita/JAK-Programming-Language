package token

type TokenType string
type Token struct {
	Type    TokenType
	Literal string
}

const (
	TT_ILLEGAL = "ILLEGAL"
	TT_EOF     = "EOF"

	TT_IDENT = "IDENT"
	TT_INT   = "INT"

	TT_ASSIGN = "="
	TT_PLUS   = "+"
	TT_MINUS  = "-"
	TT_BANG   = "!"
	TT_ASTER  = "*"
	TT_SLASH  = "/"

	TT_COMMA     = ","
	TT_SEMICOLON = ";"
	TT_LPAREN    = "("
	TT_RPAREN    = ")"
	TT_LBRACE    = "{"
	TT_RBRACE    = "}"

	TT_FUNCTION = "FUNCTION"
	TT_LET      = "LET"
)