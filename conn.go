package main

import (
	"github.com/jhump/protoreflect/grpcreflect"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	rpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

func NewGRPCConnection(ctx context.Context, addr string, insecure bool) (*grpc.ClientConn, error) {
	var dialOpts []grpc.DialOption
	if insecure {
		dialOpts = append(dialOpts, grpc.WithInsecure())
	} else {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
	}

	return grpc.DialContext(ctx, addr, dialOpts...)
}

func NewServerReflectionClient(ctx context.Context, conn *grpc.ClientConn) *grpcreflect.Client {
	cli := rpb.NewServerReflectionClient(conn)
	return grpcreflect.NewClient(ctx, cli)
}
