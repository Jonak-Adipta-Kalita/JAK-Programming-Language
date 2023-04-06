package main

import "os"

func main() {
	fileName := os.Args[1]
	file, err := os.Open(fileName)

	if err != nil {
		panic(err)
	}

	defer file.Close()
}