package main

const (
	TT_EOF = iota
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

type Token int

var tokens = []string{
	TT_EOF:         "EOF",
	TT_ILLEGAL:     "ILLEGAL",
	TT_WHITE_SPACE: "WHITE_SPACE",
	TT_INT:         "INT",
	TT_FLOAT:       "FLOAT",
	TT_STRING:      "STRING",
	TT_BOOL:        "BOOL",
	TT_PLUS:        "PLUS",
	TT_MINUS:       "MINUS",
	TT_MUL:         "MUL",
	TT_DIV:         "DIV",
	TT_POW:         "POW",
	TT_ASSIGN:      "ASSIGN",
	TT_LPAREN:      "LPAREN",
	TT_RPAREN:      "RPAREN",
	TT_LSQUARE:     "LSQUARE",
	TT_RSQUARE:     "RSQUARE",
}

func (t Token) String() string {
	return tokens[t]
}