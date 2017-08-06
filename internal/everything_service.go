package test

import (
	pb "github.com/kazegusuri/grpcurl/testdata"
	"golang.org/x/net/context"
)

type EverythingService struct{}

func NewEverythingService() *EverythingService {
	return &EverythingService{}
}

func (s *EverythingService) Simple(ctx context.Context, in *pb.SimpleMessage) (*pb.SimpleMessage, error) {
	return in, nil
}

func (s *EverythingService) Number(ctx context.Context, in *pb.NumberMessage) (*pb.NumberMessage, error) {
	return in, nil
}

func (s *EverythingService) Enum(ctx context.Context, in *pb.EnumMessage) (*pb.EnumMessage, error) {
	return in, nil
}

func (s *EverythingService) Oneof(ctx context.Context, in *pb.OneofMessage) (*pb.OneofMessage, error) {
	return in, nil
}

func (s *EverythingService) Map(ctx context.Context, in *pb.MapMessage) (*pb.MapMessage, error) {
	return in, nil
}
