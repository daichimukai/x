package main

import (
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/daichimukai/x/syakyo/monkey/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		log.Fatalf("failed to get user: %v", err)
	}
	fmt.Printf("Hello %s! This is the Monkey programming language!\n", user.Username)
	repl.Start(os.Stdin, os.Stdout)

}
