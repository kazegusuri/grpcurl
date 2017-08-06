package test

import (
	"io"

	pb "github.com/kazegusuri/grpcurl/internal/testdata"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type EchoService struct{}

func NewEchoService() *EchoService {
	return &EchoService{}
}

func (s *EchoService) Echo(ctx context.Context, in *pb.EchoMessage) (*pb.EchoMessage, error) {
	if in.ErrorCode != 0 {
		return nil, grpc.Errorf(codes.Code(in.ErrorCode), "error msg: %v", in.Value)
	}

	return &pb.EchoMessage{Value: in.Value}, nil
}

func (s *EchoService) ClientStreamingEcho(stream pb.Echo_ClientStreamingEchoServer) error {
	value := ""
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		value = msg.Value
		if msg.ErrorCode != 0 {
			return grpc.Errorf(codes.Code(msg.ErrorCode), "error msg: %v", msg.Value)
		}
	}

	if err := stream.SendAndClose(&pb.EchoMessage{Value: value}); err != nil {
		return err
	}

	return nil
}

func (s *EchoService) ServerStreamingEcho(in *pb.EchoMessage, stream pb.Echo_ServerStreamingEchoServer) error {
	if in.ErrorCode != 0 {
		return grpc.Errorf(codes.Code(in.ErrorCode), "error msg: %v", in.Value)
	}

	for i := 0; i < 10; i++ {
		if err := stream.Send(&pb.EchoMessage{Value: in.Value}); err != nil {
			return err
		}
	}
	return nil
}

func (s *EchoService) BidiStreamingBulkEcho(stream pb.Echo_BidiStreamingBulkEchoServer) error {
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if msg.ErrorCode != 0 {
			return grpc.Errorf(codes.Code(msg.ErrorCode), "error msg: %v", msg.Value)
		}

		if err := stream.Send(msg); err != nil {
			return err
		}
	}

	return nil
}
