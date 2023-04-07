package token

type TokenType string
type Token struct {
	Type    TokenType
	Literal string
}

const (
	TT_ILLEGAL = "ILLEGAL"
	TT_EOF     = "EOF"

	TT_IDENTIFIER = "IDENTIFIER"
	TT_INT        = "INT"
	TT_STRING     = "STRING"
	TT_BOOL       = "BOOL"

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
	TT_VAR      = "VAR"
	TT_PRINT    = "PRINT"
)

var keywords = map[string]TokenType{
	"fn":    TT_FUNCTION,
	"var":   TT_VAR,
	"print": TT_PRINT,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return TT_IDENTIFIER
}
