package main

import (
	"app/app"
	"flag"
	"fmt"
	"os"
)

func main() {

	var dir string
	var debug bool

	flag.StringVar(&dir, "dir", ".", "markdown directory defautl to current folder")
	flag.BoolVar(&debug, "debug", false, "do you want to print results")
	// Define a command-line flag for the rune argument
	flag.Parse()

	if err := app.Run(dir, debug); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
