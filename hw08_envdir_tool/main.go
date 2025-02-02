package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("not enough arguments")
		os.Exit(1)
	}

	env, err := ReadDir(os.Args[1])
	if err != nil {
		fmt.Printf("Error reading environment directory %s: %v", os.Args[1], err)
		os.Exit(1)
	}

	os.Exit(RunCmd(os.Args[2:], env))
}
