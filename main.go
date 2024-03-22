/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import "github.com/erikwj/brokenlinks/cmd"

// Gets set at build time via `-ldflags "-X main.sha=<value>"`

func main() {
	cmd.Execute()
}
