package ast

import "github.com/Jonak-Adipta-Kalita/JAK-Programming-Language/token"

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }