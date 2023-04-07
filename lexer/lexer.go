package lexer

func isWhiteSpace(char rune) bool {
	return char == ' ' || char == '\t' || char == '\n'
}

func isDigit(char rune) bool {
	return '0' <= char && char <= '9'
}

func isLetter(char rune) bool {
	return ('a' <= char && char <= 'z') || ('A' <= char && char <= 'Z') || char == '_'
}
