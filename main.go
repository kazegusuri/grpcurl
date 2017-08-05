package main

import (
	"os"
)

func main() {
	if err := NewRootCommand(os.Stdin).Command().Execute(); err != nil {
		os.Exit(1)
	}
}
