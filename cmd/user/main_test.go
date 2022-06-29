package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"mynote/etcdserver"
	pb "mynote/idl/userdemo"
	"mynote/pkg/constants"
	"net"
	"testing"
)

var useraddr = "127.0.0.1:8892"

func TestName(t *testing.T) {
	addr, err := net.Listen("tcp", useraddr)
	if err != nil {
		panic(err)
	}
	Init()
	server := grpc.NewServer()

	pb.RegisterUserServiceServer(server, new(UserServiceImpl))
	reflection.Register(server)

	reg, err := etcdserver.NewService(etcdserver.ServiceInfo{

		Name: constants.UserServiceName,
		IP:   useraddr, //grpc服务节点ip
	}, []string{constants.EtcdAddress}) // etcd的节点ip
	if err != nil {
		log.Fatal(err)
	}
	go reg.Start()

	err = server.Serve(addr)
	if err != nil {
		panic(err)
	}
}
