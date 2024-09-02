package usecase

import (
	"context"

	"github.com/iamrajiv/helloworld-grpc-gateway/proto/api"
)

type Server struct {
	api.UnimplementedGreeterServer
}

func (*Server) SayHello(ctx context.Context, in *api.HelloRequest) (*api.HelloReply, error) {
	//logic
	return &api.HelloReply{
		Message: "hi",
	}, nil
}

// Adds two numbers
func (*Server) AddNumbers(ctx context.Context, in *api.AddNumbersRequest) (*api.AddNumbersResponse, error) {
	return &api.AddNumbersResponse{
		Result: in.GetA() + in.GetB(),
	}, nil
}
