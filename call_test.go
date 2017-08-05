package main

import (
	"os"
	"strings"
)

func ExampleCall() {
	cmd := NewRootCommand(strings.NewReader(`{"value": "hello"}`), os.Stdout)
	cmd.Command().SetArgs([]string{"-k", "call", addr, "grpcurl.test.Echo.Echo"})
	cmd.Command().Execute()
	// Output:
	// {"value":"hello","error_code":0}
}
