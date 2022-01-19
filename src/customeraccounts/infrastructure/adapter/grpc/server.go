package customergrpc

import (
	"context"

	pb "github.com/alikarimii/go_starter/src/customeraccounts/infrastructure/adapter/grpc/proto"
)

func NewCustomerServer() pb.CustomerServer {
	return &customerServer{}
}

type customerServer struct {
}

func (c *customerServer) SignIn(ctx context.Context, in *pb.SignInReq) (*pb.SignInRes, error) {
	return &pb.SignInRes{}, nil
}
