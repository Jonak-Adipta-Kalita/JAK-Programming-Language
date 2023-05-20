package token

type TokenType string
type Token struct {
	Type     TokenType
	Literal  string
	Line     int
	PosStart int
	PosEnd   int
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENTIFIER = "IDENTIFIER"
	INT        = "INT"
	FLOAT      = "FLOAT"
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
	DOT       = "."
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
	FOREACH  = "FOREACH"
	IN       = "IN"
	IMPORT   = "USE"
	SWITCH   = "SWITCH"
	CASE     = "CASE"
	DEFAULT  = "DEFAULT"
	MACRO    = "MACRO"
)

var keywords = map[string]TokenType{
	"func":    FUNCTION,
	"var":     VAR,
	"mut":     MUTATE,
	"if":      IF,
	"elif":    ELIF,
	"else":    ELSE,
	"return":  RETURN,
	"true":    TRUE,
	"false":   FALSE,
	"null":    NULL,
	"for":     FOR,
	"foreach": FOREACH,
	"in":      IN,
	"use":     IMPORT,
	"switch":  SWITCH,
	"case":    CASE,
	"default": DEFAULT,
	"macro":   MACRO,
}

func LookupIdentifier(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENTIFIER
}
