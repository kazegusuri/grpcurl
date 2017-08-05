package main

import (
	"os"
)

func main() {
	if err := NewRootCommand().Command().Execute(); err != nil {
		os.Exit(1)
	}
}
