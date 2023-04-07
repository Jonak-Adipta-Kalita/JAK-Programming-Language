package main

import (
	"os"

	"github.com/Jonak-Adipta-Kalita/JAK-Programming-Language/repl"
)

func main() {
	if len(os.Args) != 2 {
		repl.Start(os.Stdin, os.Stdout)
	} else {
		fileName := os.Args[1]

		file, err := os.Open(fileName)

		if err != nil {
			panic(err)
		}

		defer file.Close()

		// TODO: Run Code in file
	}
}
