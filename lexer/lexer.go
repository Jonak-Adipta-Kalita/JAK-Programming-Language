package lexer

import (
	"strings"

	"github.com/Jonak-Adipta-Kalita/JAK-Programming-Language/token"
)

var escapeCharacters = map[byte]byte{
	'n':  '\n',
	't':  '\t',
	'r':  '\r',
	'"':  '"',
	'\\': '\\',
}

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
	line         int
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '#':
		l.skipSingleLineComment()
		return l.NextToken()
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = newToken(token.EQ, l.line, string(ch)+string(l.ch), l.position)
		} else {
			tok = newToken(token.ASSIGN, l.line, string(l.ch), l.position)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.line, string(l.ch), l.position)
	case '(':
		tok = newToken(token.LPAREN, l.line, string(l.ch), l.position)
	case ')':
		tok = newToken(token.RPAREN, l.line, string(l.ch), l.position)
	case ',':
		tok = newToken(token.COMMA, l.line, string(l.ch), l.position)
	case '.':
		tok = newToken(token.DOT, l.line, string(l.ch), l.position)
	case '+':
		if l.peekChar() == '+' {
			ch := l.ch
			l.readChar()
			tok = newToken(token.PLUS_PLUS, l.line, string(ch)+string(l.ch), l.position)
		} else {
			tok = newToken(token.PLUS, l.line, string(l.ch), l.position)
		}
	case '-':
		if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			tok = newToken(token.MINUS_MINUS, l.line, string(ch)+string(l.ch), l.position)
		} else {
			tok = newToken(token.MINUS, l.line, string(l.ch), l.position)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = newToken(token.NOT_EQ, l.line, string(ch)+string(l.ch), l.position)
		} else {
			tok = newToken(token.BANG, l.line, string(l.ch), l.position)
		}
	case '*':
		tok = newToken(token.ASTERISK, l.line, string(l.ch), l.position)
	case '/':
		tok = newToken(token.SLASH, l.line, string(l.ch), l.position)
	case '%':
		tok = newToken(token.MODULO, l.line, string(l.ch), l.position)
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = newToken(token.AND, l.line, string(ch)+string(l.ch), l.position)
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = newToken(token.OR, l.line, string(ch)+string(l.ch), l.position)
		}
	case '<':
		tok = newToken(token.LT, l.line, string(l.ch), l.position)
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = newToken(token.LT_EQ, l.line, string(ch)+string(l.ch), l.position)

		}
	case '>':
		tok = newToken(token.GT, l.line, string(l.ch), l.position)
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = newToken(token.GT_EQ, l.line, string(ch)+string(l.ch), l.position)
		}
	case '{':
		tok = newToken(token.LBRACE, l.line, string(l.ch), l.position)
	case '}':
		tok = newToken(token.RBRACE, l.line, string(l.ch), l.position)
	case '[':
		tok = newToken(token.LBRACKET, l.line, string(l.ch), l.position)
	case ']':
		tok = newToken(token.RBRACKET, l.line, string(l.ch), l.position)
	case '"':
		tok = newToken(token.STRING, l.line, l.readString(), l.position)
	case ':':
		tok = newToken(token.COLON, l.line, string(l.ch), l.position)
	case '^':
		tok = newToken(token.CARET, l.line, string(l.ch), l.position)
	case 0:
		tok = newToken(token.EOF, l.line, string(l.ch), l.position)
	default:
		if isLetter(l.ch) {
			indentifier := l.readIdentifier()
			tok = newToken(token.LookupIdentifier(indentifier), l.line, indentifier, l.position)

			return tok
		} else if isDigit(l.ch) {
			return l.readDecimal()
		} else {
			tok = newToken(token.ILLEGAL, l.line, string(l.ch), l.position)
		}
	}

	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, line int, literal string, posStart int) token.Token {
	return token.Token{Type: tokenType, Literal: literal, Line: line, PosStart: posStart, PosEnd: posStart + len(literal)}
}

func (l *Lexer) readIdentifier() string {
	id := ""

	position := l.position
	rposition := l.readPosition
	for isLetter(l.ch) {
		id += string(l.ch)
		l.readChar()
	}

	if strings.Contains(id, ".") {
		if !strings.HasPrefix(id, "directory.") &&
			!strings.HasPrefix(id, "file.") &&
			!strings.HasPrefix(id, "math.") &&
			!strings.HasPrefix(id, "os.") &&
			!strings.HasPrefix(id, "string.") {

			offset := strings.Index(id, ".")
			id = id[:offset]

			l.position = position
			l.readPosition = rposition
			for offset > 0 {
				l.readChar()
				offset--
			}
		}
	}

	return id
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		if l.ch == '\n' {
			l.line++
		}
		l.readChar()
	}
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	l.position++
	var escaped strings.Builder
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}

		if l.ch == '\\' {
			l.readChar()

			if escapedChar, ok := escapeCharacters[byte(l.ch)]; ok {
				escaped.WriteByte(escapedChar)
			} else {
				escaped.WriteByte('\\')
				escaped.WriteByte(byte(l.ch))
			}
		} else {
			escaped.WriteByte(byte(l.ch))
		}
	}

	return escaped.String()
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isOperator(ch byte) bool {
	return ch == byte('+') || ch == byte('%') || ch == byte('-') || ch == byte('/') || ch == byte('*')
}

func isComparison(ch byte) bool {
	return ch == byte('=') || ch == byte('!') || ch == byte('>') || ch == byte('<')
}

func isCompound(ch byte) bool {
	return ch == byte(',') || ch == byte(':') || ch == byte('"') || ch == byte(';')
}

func isBrace(ch byte) bool {
	return ch == byte('{') || ch == byte('}')
}

func isBracket(ch byte) bool {
	return ch == byte('[') || ch == byte(']')
}

func isParen(ch byte) bool {
	return ch == byte('(') || ch == byte(')')
}

func isEmpty(ch byte) bool {
	return byte(0) == ch
}

func isWhitespace(ch byte) bool {
	return ch == byte(' ') || ch == byte('\t') || ch == byte('\n') || ch == byte('\r')
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) skipSingleLineComment() {
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
	l.skipWhitespace()
}

func (l *Lexer) readDecimal() token.Token {
	integer := l.readNumber()
	if l.ch == '.' {
		l.readChar()
		fraction := l.readNumber()
		if isEmpty(l.ch) || isWhitespace(l.ch) || isOperator(l.ch) || isComparison(l.ch) || isCompound(l.ch) || isBracket(l.ch) || isBrace(l.ch) || isParen(l.ch) {
			return newToken(token.FLOAT, l.line, integer+"."+fraction, l.position)
		}
		illegalPart := l.readUntilWhitespace()
		return newToken(token.ILLEGAL, l.line, integer+"."+fraction+illegalPart, l.position)

	} else if isEmpty(l.ch) || isWhitespace(l.ch) || isOperator(l.ch) || isComparison(l.ch) || isCompound(l.ch) || isBracket(l.ch) || isBrace(l.ch) || isParen(l.ch) {
		return newToken(token.INT, l.line, integer, l.position)
	} else {
		illegalPart := l.readUntilWhitespace()
		return newToken(token.ILLEGAL, l.line, integer+illegalPart, l.position)
	}
}

func (l *Lexer) readUntilWhitespace() string {
	position := l.position
	for !isWhitespace(l.ch) {
		l.readChar()
	}
	return string(l.input[position:l.position])
}
