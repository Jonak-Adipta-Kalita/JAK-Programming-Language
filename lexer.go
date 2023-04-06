package main

func isWhiteSpace(char byte) bool {
	return char == ' ' || char == '\t' || char == '\n'
}

func isDigit(char byte) bool {
	return char >= '0' && char <= '9'
}

func isLetter(char byte) bool {
	return char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' || char == '_'
}

var eof = rune(0)
