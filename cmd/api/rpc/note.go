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

import "C"
import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"log"
	"mynote/etcdserver"
	"mynote/idl/notedemo"
	pb "mynote/idl/notedemo"
	"mynote/pkg/constants"
	"mynote/pkg/errno"
)

var noteClient notedemo.NoteServiceClient

func initNoteRpc() {
	r := etcdserver.NewResolver([]string{
		constants.EtcdAddress,
	}, constants.NoteServiceName)
	resolver.Register(r)
	conn, err := grpc.Dial(r.Scheme()+"://"+"/", grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(constants.GrpcServiceConfig),
		grpc.WithInsecure(),
	)

	//conn, err := grpc.Dial(":8972", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}

	c := pb.NewNoteServiceClient(conn)
	noteClient = c
}

// CreateNote create note info
func CreateNote(ctx context.Context, req *notedemo.CreateNoteRequest) error {
	resp, err := noteClient.CreateNote(ctx, req)
	if err != nil {
		return err
	}
	if resp.BaseResp.StatusCode != 0 {
		return errno.NewErrNo(resp.BaseResp.StatusCode, resp.BaseResp.StatusMessage)
	}
	return nil
}

// QueryNotes query list of note info
func QueryNotes(ctx context.Context, req *notedemo.QueryNoteRequest) ([]*notedemo.Note, int64, error) {
	resp, err := noteClient.QueryNote(ctx, req)
	if err != nil {
		return nil, 0, err
	}
	if resp.BaseResp.StatusCode != 0 {
		return nil, 0, errno.NewErrNo(resp.BaseResp.StatusCode, resp.BaseResp.StatusMessage)
	}
	return resp.Notes, resp.Total, nil
}

// UpdateNote update note info
func UpdateNote(ctx context.Context, req *notedemo.UpdateNoteRequest) error {
	resp, err := noteClient.UpdateNote(ctx, req)
	if err != nil {
		return err
	}
	if resp.BaseResp.StatusCode != 0 {
		return errno.NewErrNo(resp.BaseResp.StatusCode, resp.BaseResp.StatusMessage)
	}
	return nil
}

// DeleteNote delete note info
func DeleteNote(ctx context.Context, req *notedemo.DeleteNoteRequest) error {
	resp, err := noteClient.DeleteNote(ctx, req)
	if err != nil {
		return err
	}
	if resp.BaseResp.StatusCode != 0 {
		return errno.NewErrNo(resp.BaseResp.StatusCode, resp.BaseResp.StatusMessage)
	}
	return nil
}
