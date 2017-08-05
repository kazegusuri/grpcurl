package main

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/kazegusuri/grpcurl/test"
)

var (
	testPort = 35982
	addr     = fmt.Sprintf("localhost:%d", testPort)
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	go func() {
		if err := test.RunServer(ctx, testPort); err != nil {
			fmt.Fprintf(os.Stderr, "failed to start server: %v", err)
			os.Exit(1)
		}
	}()
	os.Exit(m.Run())
}
