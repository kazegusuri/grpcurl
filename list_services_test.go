package main

import (
	"os"
	"strings"
)

func ExampleListServices() {
	cmd := NewRootCommand(strings.NewReader(""), os.Stdout)
	cmd.Command().SetArgs([]string{"-k", "list_services", addr})
	cmd.Command().Execute()
	// Unordered Output:
	// grpc.reflection.v1alpha.ServerReflection
	// grpcurl.test.Echo
	// grpcurl.test.Everything
}

func ExampleListServicesMethod() {
	cmd := NewRootCommand(strings.NewReader(""), os.Stdout)
	cmd.Command().SetArgs([]string{"-k", "list_services", addr, "grpcurl.test.Echo"})
	cmd.Command().Execute()
	// Unordered Output:
	// grpcurl.test.Echo.Echo
	// grpcurl.test.Echo.ClientStreamingEcho
	// grpcurl.test.Echo.ServerStreamingEcho
	// grpcurl.test.Echo.BidiStreamingBulkEcho
}

func ExampleListServicesMethodLong() {
	cmd := NewRootCommand(strings.NewReader(""), os.Stdout)
	cmd.Command().SetArgs([]string{"-k", "list_services", addr, "-l", "grpcurl.test.Echo"})
	cmd.Command().Execute()
	// Unordered Output:
	// grpcurl.test.Echo.Echo(grpcurl.test.EchoMessage) return (grpcurl.test.EchoMessage)
	// grpcurl.test.Echo.ClientStreamingEcho(streaming grpcurl.test.EchoMessage) return (grpcurl.test.EchoMessage)
	// grpcurl.test.Echo.ServerStreamingEcho(grpcurl.test.EchoMessage) return (streaming grpcurl.test.EchoMessage)
	// grpcurl.test.Echo.BidiStreamingBulkEcho(streaming grpcurl.test.EchoMessage) return (streaming grpcurl.test.EchoMessage)
}
