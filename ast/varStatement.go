package ast

import "github.com/Jonak-Adipta-Kalita/JAK-Programming-Language/token"

type VarStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *VarStatement) statementNode() {}

func (ls *VarStatement) TokenLiteral() string { return ls.Token.Literal }
