package main

import (
	"fmt"
	"os"

	"github.com/wawandco/facto"
)

func main() {
	run(os.Args)
}

func run(args []string) {
	wd, err := os.Getwd()
	if err != nil {
		handleErr(err)
	}

	if len(args) < 2 {
		info()
		return
	}

	args = args[1:]
	if args[0] != "generate" {
		info()
		return
	}

	err = facto.Generate(wd, args)
	if err != nil {
		handleErr(err)
	}
}

// A simple info command to guide users.
func info() {
	fmt.Print("Usage:\n facto <command>\n\n")
	fmt.Println("Valid commands are `generate` and `info`.")
}

func handleErr(err error) {
	fmt.Printf("Error running command:\n %v\n\n", err)
	info()
	os.Exit(1)
}
