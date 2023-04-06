package main

import "fmt"

type Token struct {
	Type  string
	Value interface{}
}

func NewToken(type_ string, value interface{}) *Token {
	return &Token{Type: type_, Value: value}
}

func (t *Token) String() string {
	if t.Value != nil {
		return fmt.Sprintf("%v:%v", t.Type, t.Value)
	}
	return t.Type
}
