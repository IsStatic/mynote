package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"mynote/cmd/note/dal"
	"mynote/cmd/note/rpc"
	"mynote/etcdserver"
	pb "mynote/idl/notedemo"
	"mynote/pkg/constants"
	"mynote/pkg/tracer"
	"net"
)

func Init() {
	tracer.InitJaeger(constants.NoteServiceName)
	rpc.InitRPC()
	dal.Init()
}

func main() {
	addr, err := net.Listen("tcp", constants.NoteAddress)
	if err != nil {
		panic(err)
	}
	Init()
	server := grpc.NewServer()

	pb.RegisterNoteServiceServer(server, new(NoteServiceImpl))
	reflection.Register(server)

	reg, err := etcdserver.NewService(etcdserver.ServiceInfo{

		Name: constants.NoteServiceName,
		IP:   constants.NoteAddress, //grpc服务节点ip
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
