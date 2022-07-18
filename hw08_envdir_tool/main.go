package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Parse()
	parameters := flag.Args()

	if len(parameters) < 2 {
		fmt.Println(
			"Error occurred: Invalid input params.\nCommand required at least 2 parameters:",
			"\nAt first, directory with env parameters.\nAt second, command to execute.",
			"\nAnd next, parameters what will be passed to command afterwards",
		)
		os.Exit(1)
	}

	environment, _ := ReadDir(parameters[0])

	RunCmd(parameters[1:], environment)
}
