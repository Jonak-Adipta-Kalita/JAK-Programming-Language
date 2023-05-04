package token

type TokenType string
type Token struct {
	Type    TokenType
	Literal string
	Line    int
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENTIFIER = "IDENTIFIER"
	INT        = "INT"
	STRING     = "STRING"
	NULL       = "NULL"

	ASSIGN      = "="
	PLUS        = "+"
	MINUS       = "-"
	BANG        = "!"
	ASTERISK    = "*"
	SLASH       = "/"
	MODULO      = "%"
	LT          = "<"
	GT          = ">"
	EQ          = "=="
	NOT_EQ      = "!="
	LT_EQ       = "<="
	GT_EQ       = ">="
	PLUS_PLUS   = "++"
	MINUS_MINUS = "--"
	AND         = "&&"
	OR          = "||"
	CARET       = "^"

	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	LBRACKET  = "["
	RBRACKET  = "]"

	FUNCTION = "FUNCTION"
	VAR      = "VAR"
	MUTATE   = "MUTATE"
	IF       = "IF"
	ELIF     = "ELIF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	FOR      = "FOR"
	IMPORT   = "USE"
)

var keywords = map[string]TokenType{
	"func":   FUNCTION,
	"var":    VAR,
	"mut":    MUTATE,
	"if":     IF,
	"elif":   ELIF,
	"else":   ELSE,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
	"null":   NULL,
	"for":    FOR,
	"use":    IMPORT,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENTIFIER
}
