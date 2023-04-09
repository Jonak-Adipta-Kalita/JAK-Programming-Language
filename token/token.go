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
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	FOR      = "FOR"
	IMPORT   = "USE"
	AS       = "AS"
)

var keywords = map[string]TokenType{
	"func":   FUNCTION,
	"var":    VAR,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
	"for":    FOR,
	"use":    IMPORT,
	"as":     AS,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENTIFIER
}
