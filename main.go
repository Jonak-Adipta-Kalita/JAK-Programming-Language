package main

import (
	"fmt"
	"os"
)

func main() {
	fileName := os.Args[1]
	
	if fileName == "" {
		for {
			var codeLine string
			fmt.Print(">>> ")
			fmt.Scanln(&codeLine)

			// TODO: Run code
		}
	}

	file, err := os.Open(fileName)

	if err != nil {
		panic(err)
	}

	defer file.Close()
}