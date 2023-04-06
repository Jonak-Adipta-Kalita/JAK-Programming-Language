package main

type Token int

var tokens = []string{
	TT_EOF:     "EOF",
	TT_ILLEGAL: "ILLEGAL",
	TT_INT:     "INT",
	TT_FLOAT:   "FLOAT",
	TT_STRING:  "STRING",
	TT_BOOL:    "BOOL",
	TT_PLUS:    "PLUS",
	TT_MINUS:   "MINUS",
	TT_MUL:     "MUL",
	TT_DIV:     "DIV",
	TT_POW:     "POW",
	TT_ASSIGN:  "ASSIGN",
	TT_LPAREN:  "LPAREN",
	TT_RPAREN:  "RPAREN",
	TT_LSQUARE: "LSQUARE",
	TT_RSQUARE: "RSQUARE",
}

func (t Token) String() string {
	return tokens[t]
}