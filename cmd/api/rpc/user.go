// Copyright 2021 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

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
		grpc.WithDefaultServiceConfig(constants.GrpcServiceConfig),
		grpc.WithInsecure(),
	)

	//conn, err := grpc.Dial(":8972", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}

	c := pb.NewUserServiceClient(conn)
	userClient = c

}

// CreateUser create user info
func CreateUser(ctx context.Context, req *userdemo.CreateUserRequest) error {
	resp, err := userClient.CreateUser(ctx, req)
	if err != nil {
		return err
	}
	if resp.BaseResp.StatusCode != 0 {
		return errno.NewErrNo(resp.BaseResp.StatusCode, resp.BaseResp.StatusMessage)
	}
	return nil
}

// CheckUser check user info
func CheckUser(ctx context.Context, req *userdemo.CheckUserRequest) (int64, error) {
	resp, err := userClient.CheckUser(ctx, req)
	if err != nil {
		return 0, err
	}
	if resp.BaseResp.StatusCode != 0 {
		return 0, errno.NewErrNo(resp.BaseResp.StatusCode, resp.BaseResp.StatusMessage)
	}
	return resp.UserId, nil
}
