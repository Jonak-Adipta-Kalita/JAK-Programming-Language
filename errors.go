package main

import "fmt"

type Error struct {
    posStart   *Position
    posEnd     *Position
    errorName  string
    details    string
}

func NewError(posStart, posEnd *Position, errorName, details string) *Error {
    return &Error{posStart, posEnd, errorName, details}
}

func (e *Error) AsString() string {
    result := fmt.Sprintf("%s: %s\n", e.errorName, e.details)
    result += fmt.Sprintf("File %s, line %d\n", e.posStart.Fn, e.posStart.Ln + 1)
    return result
}

type IllegalCharError struct {
    *Error
}

func NewIllegalCharError(posStart, posEnd *Position, details string) *IllegalCharError {
    return &IllegalCharError{NewError(posStart, posEnd, "Illegal Character", details)}
}