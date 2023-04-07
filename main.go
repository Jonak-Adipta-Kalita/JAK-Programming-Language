package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		for {
			var codeLine string
			fmt.Print(">>> ")
			fmt.Scanln(&codeLine)
			fmt.Println(codeLine)

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