package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		reader := bufio.NewReader(os.Stdin)
		for {
			var codeLine string
			fmt.Print(">>> ")
			codeLine, _ = reader.ReadString('\n')
			fmt.Print(codeLine + "\n")

			// TODO: Run code
		}
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