package test

import (
	pb "github.com/kazegusuri/grpcurl/testdata"
	"golang.org/x/net/context"
)

type EchoServiceV2 struct{}

func NewEchoServiceV2() *EchoServiceV2 {
	return &EchoServiceV2{}
}

func (s *EchoServiceV2) Echo(ctx context.Context, in *pb.EchoMessage) (*pb.EchoMessage, error) {
	return &pb.EchoMessage{Value: in.Value}, nil
}
