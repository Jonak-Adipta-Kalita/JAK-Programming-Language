package token

type TokenType string
type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENTIFIER = "IDENTIFIER"
	INT        = "INT"
	STRING     = "STRING"
	BOOL       = "BOOL"

	ASSIGN = "="
	PLUS   = "+"
	MINUS  = "-"
	BANG   = "!"
	ASTER  = "*"
	SLASH  = "/"

	LT = "<"
	GT = ">"

	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"

	FUNCTION = "FUNCTION"
	VAR      = "VAR"
	PRINT    = "PRINT"
)

var keywords = map[string]TokenType{
	"fn":    FUNCTION,
	"var":   VAR,
	"print": PRINT,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENTIFIER
}
