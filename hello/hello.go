package main

import (
	"fmt"
	"toombs/greetings"
)

func main() {
	// Get a greeting message and print it.
	message := greetings.Hello("Mr. Pants")
	fmt.Println(message)
}
