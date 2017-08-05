package main

import (
	"strings"
)

func ExampleCall() {
	cmd := NewRootCommand(strings.NewReader(`{"value": "hello"}`))
	cmd.Command().SetArgs([]string{"-k", "call", addr, "grpcurl.test.Echo.Echo"})
	cmd.Command().Execute()
	// Output:
	// {"value":"hello","error_code":0}
}
