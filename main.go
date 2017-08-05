package main

import (
	"os"
)

func main() {
	if err := NewRootCommand(os.Stdin, os.Stdout).Command().Execute(); err != nil {
		os.Exit(1)
	}
}
