package test

import (
	"context"
	"fmt"
	"net"

	pb "github.com/kazegusuri/grpcurl/testdata"
	pbv2 "github.com/kazegusuri/grpcurl/testdata/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func RunServer(ctx context.Context, port int) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return fmt.Errorf("failed to list: %v", err)
	}
	s := grpc.NewServer()
	defer s.Stop()

	go func() {
		pb.RegisterEchoServer(s, NewEchoService())
		pbv2.RegisterEchoServer(s, NewEchoServiceV2())
		pb.RegisterEverythingServer(s, NewEverythingService())
		reflection.Register(s)

		s.Serve(l)
		cancel()
	}()

	select {
	case <-ctx.Done():
		s.Stop()
	}

	return nil
}
