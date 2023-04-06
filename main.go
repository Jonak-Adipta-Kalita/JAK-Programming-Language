package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(">>> ")
	text, _ := reader.ReadString('\n')
	fmt.Print(text)
}

// lexer
// 1. read the input
// 2. break it into tokens
// 3. return the tokens one at a time

// parser
// 1. read the tokens
// 2. determine if the syntax is valid
// 3. if valid, return the AST
// 4. if invalid, return an error

// AST
// 1. the AST is just a tree of nodes
// 2. each node is an expression
// 3. each expression has a type
// 4. each expression has zero or more children
