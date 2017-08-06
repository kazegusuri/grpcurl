package main

import (
	"context"
	"flag"

	"github.com/golang/glog"
	"github.com/kazegusuri/grpcurl/test"
)

var (
	port = flag.Int("port", 8888, "port")
)

func main() {
	flag.Parse()
	defer glog.Flush()

	ctx := context.Background()
	if err := test.RunServer(ctx, *port); err != nil {
		glog.Exit(err)
	}
}
