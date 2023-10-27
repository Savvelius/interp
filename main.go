package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/Savvelius/go-interp/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello, %v. This is the monkey programming language.\n", user.Username)
	fmt.Printf("Feel free to type in any commands.\n")
	repl.Start(os.Stdin, os.Stdout)
}
