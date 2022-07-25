package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"mynote/cmd/user/dal/db"
	"mynote/etcdserver"
	pb "mynote/idl/userdemo"
	"mynote/pkg/constants"
	"mynote/pkg/tracer"
	"net"
)

func Init() {
	tracer.InitJaeger(constants.UserServiceName)
	go db.Init()
}

func main() {
	addr, err := net.Listen("tcp", constants.UserAddress)
	if err != nil {
		panic(err)
	}
	Init()
	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, new(UserServiceImpl))
	reflection.Register(server)
	reg, err := etcdserver.NewService(etcdserver.ServiceInfo{
		Name:   constants.UserServiceName,
		IP:     constants.UserAddress, //grpc服务节点ip
		Weight: constants.UserWeight,
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
