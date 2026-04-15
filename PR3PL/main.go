/*Eugenio Giusepi Montilla Russo*/
/*29958321*/

package main

import (
	"fmt"
	"os"
	"os/user"
	"pr3pl/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the PR3PL programming language!\n",
		user.Username)
	fmt.Printf("Feel free to type in commands\n")

	repl.Start(os.Stdin, os.Stdout)
}
