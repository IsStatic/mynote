package rpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"log"
	"mynote/etcdserver"
	"mynote/idl/userdemo"
	pb "mynote/idl/userdemo"
	"mynote/pkg/constants"
	"mynote/pkg/errno"
)

var userClient userdemo.UserServiceClient

func initUserRpc() {
	r := etcdserver.NewResolver([]string{
		constants.EtcdAddress,
	}, constants.UserServiceName)
	resolver.Register(r)

	conn, err := grpc.Dial(r.Scheme()+"://"+"/", grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(constants.GrpcServiceConfig))

	//conn, err := grpc.Dial(":8972", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	c := pb.NewUserServiceClient(conn)
	userClient = c
}

// MGetUser multiple get list of user info
func MGetUser(ctx context.Context, req *userdemo.MGetUserRequest) (map[int64]*userdemo.User, error) {
	resp, err := userClient.MGetUser(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.BaseResp.StatusCode != 0 {
		return nil, errno.NewErrNo(resp.BaseResp.StatusCode, resp.BaseResp.StatusMessage)
	}
	res := make(map[int64]*userdemo.User)
	for _, u := range resp.Users {
		res[u.UserId] = u
	}
	return res, nil
}
